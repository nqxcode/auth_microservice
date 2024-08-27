package app

import (
	"context"
	"log"

	"github.com/nqxcode/auth_microservice/internal/api/auth"
	"github.com/nqxcode/auth_microservice/internal/config"
	"github.com/nqxcode/auth_microservice/internal/repository"
	logRepository "github.com/nqxcode/auth_microservice/internal/repository/log"
	pgUserRepository "github.com/nqxcode/auth_microservice/internal/repository/user/pg"
	redisUserRepository "github.com/nqxcode/auth_microservice/internal/repository/user/redis"
	"github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/async"
	authService "github.com/nqxcode/auth_microservice/internal/service/auth"
	cacheUserService "github.com/nqxcode/auth_microservice/internal/service/cache/user"
	hashService "github.com/nqxcode/auth_microservice/internal/service/hash"
	logService "github.com/nqxcode/auth_microservice/internal/service/log"
	"github.com/nqxcode/auth_microservice/internal/service/validator"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/nqxcode/platform_common/client/cache"
	"github.com/nqxcode/platform_common/client/cache/redis"
	"github.com/nqxcode/platform_common/client/db"
	"github.com/nqxcode/platform_common/client/db/pg"
	"github.com/nqxcode/platform_common/client/db/transaction"
	"github.com/nqxcode/platform_common/closer"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	hashingConfig config.HashingConfig
	redisConfig   cache.RedisConfig

	dbClient    db.Client
	txManager   db.TxManager
	asyncRunner async.Runner

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	userRepository repository.UserRepository
	logRepository  repository.LogRepository

	cacheUserRepository repository.UserRepository

	logService       service.LogService
	hashService      service.HashService
	authService      service.AuthService
	cacheUserService service.CacheUserService
	validatorService service.ValidatorService

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

func (s *serviceProvider) RedisConfig() cache.RedisConfig {
	if s.redisConfig == nil {
		cfg, err := config.NewRedisConfig()
		if err != nil {
			log.Fatalf("failed to get redis config: %s", err.Error())
		}

		s.redisConfig = cfg
	}

	return s.redisConfig
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

func (s *serviceProvider) RedisPool() *redigo.Pool {
	if s.redisPool == nil {
		s.redisPool = &redigo.Pool{
			MaxIdle:     s.RedisConfig().MaxIdle(),
			IdleTimeout: s.RedisConfig().IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(
					ctx,
					"tcp",
					s.RedisConfig().Address(),
					redigo.DialDatabase(s.RedisConfig().DB()),
					redigo.DialPassword(s.RedisConfig().Password()),
				)
			},
		}
	}

	return s.redisPool
}

func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

func (s *serviceProvider) AsyncRunner() async.Runner {
	if s.asyncRunner == nil {
		s.asyncRunner = async.NewRunner()
	}

	return s.asyncRunner
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = pgUserRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

func (s *serviceProvider) CacheUserRepository() repository.UserRepository {
	if s.cacheUserRepository == nil {
		s.cacheUserRepository = redisUserRepository.NewRepository(s.RedisClient())
	}

	return s.cacheUserRepository
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
			s.UserRepository(ctx),
			s.ValidatorService(),
			s.LogService(ctx),
			s.HashService(ctx),
			s.CacheUserService(),
			s.TxManager(ctx),
			s.AsyncRunner(),
		)
	}

	return s.authService
}

func (s *serviceProvider) CacheUserService() service.CacheUserService {
	if s.cacheUserService == nil {
		s.cacheUserService = cacheUserService.NewService(
			s.RedisClient(),
			s.CacheUserRepository(),
		)
	}

	return s.cacheUserService
}

func (s *serviceProvider) ValidatorService() service.ValidatorService {
	if s.validatorService == nil {
		s.validatorService = validator.NewValidator()
	}

	return s.validatorService
}

func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
