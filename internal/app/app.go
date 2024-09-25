package app

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nqxcode/auth_microservice/internal/logger"
	"github.com/nqxcode/platform_common/closer"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/nqxcode/auth_microservice/internal/config"
	"github.com/nqxcode/auth_microservice/internal/interceptor"
	desc "github.com/nqxcode/auth_microservice/pkg/auth_v1"
	_ "github.com/nqxcode/auth_microservice/pkg/statik" // nolint: revive
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

// App application
type App struct {
	serviceProvider  *serviceProvider
	grpcServer       *grpc.Server
	httpServer       *http.Server
	swaggerServer    *http.Server
	prometheusServer *http.Server
}

// NewApp new application
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// Run run application
func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)

	wg := &sync.WaitGroup{}
	errChan := make(chan error, 4)

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runGRPCServer(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to run GRPC server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runHTTPServer(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to run HTTP server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runSwaggerServer(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.serviceProvider.ConsumerService(ctx).RunConsumer(ctx)
		if err != nil {
			errChan <- fmt.Errorf("failed to run consumer: %s", err.Error())
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		err := a.runPrometheus(ctx)
		if err != nil {
			if err != nil {
				errChan <- fmt.Errorf("failed to run prometheus server: %s", err.Error())
			}
		}
	}()

	gracefulShutdown(ctx, cancel, wg)

	errs := make([]error, 0, len(errChan))
	for i := 0; i < len(errChan); i++ {
		err := <-errChan
		if err != nil {
			errs = append(errs, err)
		}
	}

	return multierr.Combine(errs...)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initMetrics,
		a.initTracing,
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	a.serviceProvider.InitLogger(ctx)
	return nil
}

func (a *App) initMetrics(ctx context.Context) error {
	a.serviceProvider.InitMetric(ctx)
	return nil
}

// initTracing initializes tracing
func (a *App) initTracing(ctx context.Context) error {
	a.serviceProvider.InitTracing(ctx)
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(configPath)
	if err != nil {
		log.Printf("No %s file found, using environment variables: %v", configPath, err)
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	tlsCert, err := tls.X509KeyPair(a.serviceProvider.GRPCConfig().Cert(), a.serviceProvider.GRPCConfig().Key())
	if err != nil {
		return err
	}

	creds := credentials.NewServerTLSFromCert(&tlsCert)

	a.grpcServer = grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(
		grpcMiddleware.ChainUnaryServer(
			interceptor.MetricsInterceptor,
			interceptor.ServerTracingInterceptor,
			interceptor.LogInterceptor,
			interceptor.ValidateInterceptor,
		),
	))

	reflection.Register(a.grpcServer)

	desc.RegisterAuthV1Server(a.grpcServer, a.serviceProvider.AuthImpl(ctx))

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterAuthV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:         a.serviceProvider.HTTPConfig().Address(),
		ReadTimeout:  a.serviceProvider.HTTPConfig().ReadTimeout(),
		WriteTimeout: a.serviceProvider.HTTPConfig().WriteTimeout(),
		IdleTimeout:  a.serviceProvider.HTTPConfig().IdleTimeout(),
		Handler:      corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:         a.serviceProvider.SwaggerConfig().Address(),
		ReadTimeout:  a.serviceProvider.SwaggerConfig().ReadTimeout(),
		WriteTimeout: a.serviceProvider.SwaggerConfig().WriteTimeout(),
		IdleTimeout:  a.serviceProvider.SwaggerConfig().IdleTimeout(),
		Handler:      mux,
	}

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info("GRPC server is running", zap.String("address", a.serviceProvider.GRPCConfig().Address()))

	go func() {
		<-ctx.Done()
		a.grpcServer.GracefulStop()
		logger.Info("GRPC server gracefully stopped")
	}()

	listener, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info("HTTP server is running", zap.String("address", a.serviceProvider.HTTPConfig().Address()))

	go func() {
		<-ctx.Done()
		if shutdownErr := a.httpServer.Shutdown(ctx); shutdownErr != nil {
			logger.Info("HTTP server shutdown err", zap.Error(shutdownErr))
			if closeErr := a.httpServer.Close(); shutdownErr != nil {
				logger.Info("HTTP server close err", zap.Error(closeErr))
			}
		} else {
			logger.Info("HTTP server gracefully stopped")
		}
	}()

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer(ctx context.Context) error {
	logger.Info("Swagger server is running", zap.String("address", a.serviceProvider.SwaggerConfig().Address()))

	go func() {
		<-ctx.Done()
		if shutdownErr := a.swaggerServer.Shutdown(ctx); shutdownErr != nil {
			logger.Info("Swagger server shutdown err", zap.Error(shutdownErr))
			if closeErr := a.swaggerServer.Close(); shutdownErr != nil {
				logger.Info("Swagger server close err", zap.Error(closeErr))
			}
		} else {
			logger.Info("Swagger server gracefully stopped")
		}
	}()

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		logger.Info("Serving swagger file", zap.String("path", path))

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Open swagger file", zap.String("path", path))

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close() // nolint: errcheck

		logger.Info("Read swagger file", zap.String("path", path))

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Write swagger file", zap.String("path", path))

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		logger.Info("Served swagger file", zap.String("path", path))
	}
}

func (a *App) runPrometheus(ctx context.Context) error {
	cfg := a.serviceProvider.NewPrometheusConfig()

	mux := http.NewServeMux()
	mux.Handle(cfg.MetricsPath(), promhttp.Handler())

	a.prometheusServer = &http.Server{
		Addr:              cfg.Address(),
		Handler:           mux,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout(),
	}

	logger.Info(fmt.Sprintf("Prometheus server is running on %s", cfg.Address()))

	go func() {
		<-ctx.Done()
		if shutdownErr := a.prometheusServer.Shutdown(ctx); shutdownErr != nil {
			logger.Info("Swagger server shutdown err", zap.Error(shutdownErr))
			if closeErr := a.prometheusServer.Close(); shutdownErr != nil {
				logger.Info("Swagger server close err", zap.Error(closeErr))
			}
		} else {
			logger.Info("Swagger server gracefully stopped")
		}
	}()

	err := a.prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-waitSignal():
		log.Println("terminating: via signal")
	}

	cancel()
	if wg != nil {
		wg.Wait()
	}
}

func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	return sigterm
}
