package app

import (
	"context"
	"log"

	"github.com/nqxcode/auth_microservice/internal/api/auth"
	"github.com/nqxcode/auth_microservice/internal/config"
	"github.com/nqxcode/auth_microservice/internal/repository"
	logRepository "github.com/nqxcode/auth_microservice/internal/repository/log"
	userRepository "github.com/nqxcode/auth_microservice/internal/repository/user"
	"github.com/nqxcode/auth_microservice/internal/service"
	authService "github.com/nqxcode/auth_microservice/internal/service/auth"
	hashService "github.com/nqxcode/auth_microservice/internal/service/hash"
	logService "github.com/nqxcode/auth_microservice/internal/service/log"
	"github.com/nqxcode/platform_common/client/db"
	"github.com/nqxcode/platform_common/client/db/pg"
	"github.com/nqxcode/platform_common/client/db/transaction"
	"github.com/nqxcode/platform_common/closer"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	hashingConfig config.HashingConfig

	dbClient       db.Client
	txManager      db.TxManager
	userRepository repository.UserRepository
	logRepository  repository.LogRepository

	logService  service.LogService
	hashService service.HashService
	authService service.AuthService

	authImpl *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) HashingConfig() config.HashingConfig {
	if s.hashingConfig == nil {
		cfg, err := config.NewHashingConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.hashingConfig = cfg
	}

	return s.hashingConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) LogService(ctx context.Context) service.LogService {
	if s.logService == nil {
		s.logService = logService.NewService(
			s.LogRepository(ctx),
		)
	}

	return s.logService
}

func (s *serviceProvider) HashService(ctx context.Context) service.HashService {
	if s.hashService == nil {
		s.hashService = hashService.NewService(
			s.HashingConfig().Salt(ctx),
		)
	}

	return s.hashService
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.AuthRepository(ctx),
			s.LogService(ctx),
			s.HashService(ctx),
			s.TxManager(ctx),
		)
	}

	return s.authService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
