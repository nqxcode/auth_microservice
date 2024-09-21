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

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/nqxcode/platform_common/closer"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/multierr"
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
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
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

	a.grpcServer = grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(interceptor.ValidateInterceptor))

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
	log.Printf("GRPC server is running on %s", a.serviceProvider.GRPCConfig().Address())

	go func() {
		<-ctx.Done()
		a.grpcServer.GracefulStop()
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
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	go func() {
		<-ctx.Done()
		if shutdownErr := a.httpServer.Shutdown(ctx); shutdownErr != nil {
			log.Printf("http server shutdown err: %s", shutdownErr)
			if closeErr := a.httpServer.Close(); shutdownErr != nil {
				log.Printf("http server close err: %s", closeErr)
			}
		}
	}()

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer(ctx context.Context) error {
	log.Printf("Swagger server is running on %s", a.serviceProvider.SwaggerConfig().Address())

	go func() {
		<-ctx.Done()
		if shutdownErr := a.swaggerServer.Shutdown(ctx); shutdownErr != nil {
			log.Printf("swagger server shutdown err: %s", shutdownErr)
			if closeErr := a.swaggerServer.Close(); shutdownErr != nil {
				log.Printf("swagger server close err: %s", closeErr)
			}
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
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close() // nolint: errcheck

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
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
