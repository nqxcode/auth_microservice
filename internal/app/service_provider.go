package app

import (
	"context"
	"log"
	"os"

	"github.com/nqxcode/auth_microservice/internal/metric"
	"github.com/nqxcode/auth_microservice/internal/tracing"

	"github.com/natefinch/lumberjack"
	"github.com/nqxcode/auth_microservice/internal/logger"
	"github.com/nqxcode/auth_microservice/internal/service/token"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/IBM/sarama"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/nqxcode/platform_common/client/broker/kafka"
	kafkaConsumer "github.com/nqxcode/platform_common/client/broker/kafka/consumer"
	"github.com/nqxcode/platform_common/client/broker/kafka/producer"
	"github.com/nqxcode/platform_common/client/cache"
	"github.com/nqxcode/platform_common/client/cache/redis"
	"github.com/nqxcode/platform_common/client/db"
	"github.com/nqxcode/platform_common/client/db/pg"
	"github.com/nqxcode/platform_common/client/db/transaction"
	"github.com/nqxcode/platform_common/closer"

	"github.com/nqxcode/auth_microservice/internal/api/auth"
	"github.com/nqxcode/auth_microservice/internal/config"
	"github.com/nqxcode/auth_microservice/internal/repository"
	accessibleRoleRepository "github.com/nqxcode/auth_microservice/internal/repository/accessible_role"
	logRepository "github.com/nqxcode/auth_microservice/internal/repository/log"
	pgUserRepository "github.com/nqxcode/auth_microservice/internal/repository/user/pg"
	redisUserRepository "github.com/nqxcode/auth_microservice/internal/repository/user/redis"
	"github.com/nqxcode/auth_microservice/internal/service"
	"github.com/nqxcode/auth_microservice/internal/service/async"
	auditLogService "github.com/nqxcode/auth_microservice/internal/service/audit_log"
	authService "github.com/nqxcode/auth_microservice/internal/service/auth"
	cacheUserService "github.com/nqxcode/auth_microservice/internal/service/cache/user"
	userSaverConsumer "github.com/nqxcode/auth_microservice/internal/service/consumer/user_saver"
	hashService "github.com/nqxcode/auth_microservice/internal/service/hash"
	auditLogSender "github.com/nqxcode/auth_microservice/internal/service/producer/audit_log_sender"
	"github.com/nqxcode/auth_microservice/internal/service/validator"
)

type serviceProvider struct {
	pgConfig            config.PGConfig
	grpcConfig          config.GRPCConfig
	httpConfig          config.HTTPConfig
	swaggerConfig       config.SwaggerConfig
	hashingConfig       config.HashingConfig
	redisConfig         cache.RedisConfig
	kafkaConsumerConfig kafka.ConsumerConfig
	kafkaProducerConfig kafka.ProducerConfig
	authConfig          config.AuthConfig
	loggerConfig        config.LoggerConfig
	prometheusConfig    config.PrometheusConfig
	appConfig           config.AppConfig

	dbClient  db.Client
	txManager db.TxManager

	redisPool   *redigo.Pool
	redisClient cache.RedisClient

	consumer             kafka.Consumer
	consumerGroup        sarama.ConsumerGroup
	consumerGroupHandler *kafkaConsumer.GroupHandler
	syncProducer         kafka.SyncProducer

	asyncRunner async.Runner

	userRepository           repository.UserRepository
	logRepository            repository.LogRepository
	cacheUserRepository      repository.UserRepository
	accessibleRoleRepository repository.AccessibleRoleRepository

	auditLogService  service.AuditLogService
	hashService      service.HashService
	authService      service.AuthService
	cacheUserService service.CacheUserService
	validatorService service.ValidatorService
	consumerService  service.ConsumerService
	producerService  service.ProducerService

	authImpl *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

// InitLogger init logger
func (s *serviceProvider) InitLogger(_ context.Context) {
	cnf := s.NewLoggerConfig()

	logLevel := cnf.GetLogLevel()
	rollingConfig := cnf.GetRollingConfig()

	stdout := zapcore.AddSync(os.Stdout)

	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	atomicLevel := zap.NewAtomicLevelAt(level)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   rollingConfig.Filename,
		MaxSize:    rollingConfig.MaxSizeInMegabytes,
		MaxBackups: rollingConfig.MaxBackups,
		MaxAge:     rollingConfig.MaxAgeInDays,
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, atomicLevel),
		zapcore.NewCore(fileEncoder, file, atomicLevel),
	)

	logger.Init(core)
}

// InitMetric init metric
func (s *serviceProvider) InitMetric(ctx context.Context) {
	err := metric.Init(ctx)
	if err != nil {
		log.Fatalf("failed to init metrics: %v", err)
	}
}

// InitTracing init tracing
func (s *serviceProvider) InitTracing(_ context.Context) {
	tracing.Init(logger.Logger(), s.NewAppConfig().GetName())
}

// PGConfig config for pg
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

// GRPCConfig config for grpc
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

// HashingConfig config for hashing
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

// RedisConfig config for redis
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

// HTTPConfig config for http server
func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := config.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}

// SwaggerConfig config for swagger
func (s *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if s.swaggerConfig == nil {
		cfg, err := config.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to get swagger config: %s", err.Error())
		}

		s.swaggerConfig = cfg
	}

	return s.swaggerConfig
}

// KafkaConsumerConfig config for kafka consumer
func (s *serviceProvider) KafkaConsumerConfig() kafka.ConsumerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := config.NewKafkaConsumerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka consumer config: %s", err.Error())
		}

		s.kafkaConsumerConfig = cfg
	}

	return s.kafkaConsumerConfig
}

// KafkaProducerConfig config for kafka producer
func (s *serviceProvider) KafkaProducerConfig() kafka.ProducerConfig {
	if s.kafkaConsumerConfig == nil {
		cfg, err := config.NewKafkaProducerConfig()
		if err != nil {
			log.Fatalf("failed to get kafka producer config: %s", err.Error())
		}

		s.kafkaProducerConfig = cfg
	}

	return s.kafkaProducerConfig
}

// NewAuthConfig config for auth
func (s *serviceProvider) NewAuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		cfg, err := config.NewAuthConfig()
		if err != nil {
			log.Fatalf("failed to get auth config: %s", err.Error())
		}

		s.authConfig = cfg
	}

	return s.authConfig
}

func (s *serviceProvider) NewAppConfig() config.AppConfig {
	if s.appConfig == nil {
		cfg, err := config.NewAppConfig()
		if err != nil {
			log.Fatalf("failed to get app config: %s", err.Error())
		}

		s.appConfig = cfg
	}

	return s.appConfig
}

// NewLoggerConfig config for logger
func (s *serviceProvider) NewLoggerConfig() config.LoggerConfig {
	if s.loggerConfig == nil {
		cfg := config.NewLoggerConfig()

		s.loggerConfig = cfg
	}

	return s.loggerConfig
}

// NewPrometheusConfig config for prometheus
func (s *serviceProvider) NewPrometheusConfig() config.PrometheusConfig {
	cfg, err := config.NewPrometheusConfig()
	if err != nil {
		log.Fatalf("failed to get prometheus config: %s", err.Error())
	}

	return cfg
}

// DBClient client for pg database
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

// RedisPool redis pool
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

// RedisClient client for redis
func (s *serviceProvider) RedisClient() cache.RedisClient {
	if s.redisClient == nil {
		s.redisClient = redis.NewClient(s.RedisPool(), s.RedisConfig())
	}

	return s.redisClient
}

// TxManager tx manager
func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}

	return s.txManager
}

// AsyncRunner runner to run handler in goroutine
func (s *serviceProvider) AsyncRunner() async.Runner {
	if s.asyncRunner == nil {
		s.asyncRunner = async.NewRunner()
	}

	return s.asyncRunner
}

// NewTokenGenerator token generator
func (s *serviceProvider) NewTokenGenerator() service.TokenGenerator {
	return token.NewGenerator()
}

// UserRepository user repository
func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = pgUserRepository.NewRepository(s.DBClient(ctx))
	}

	return s.userRepository
}

// LogRepository audit log repository
func (s *serviceProvider) LogRepository(ctx context.Context) repository.LogRepository {
	if s.logRepository == nil {
		s.logRepository = logRepository.NewRepository(s.DBClient(ctx))
	}

	return s.logRepository
}

// CacheUserRepository cache user repository
func (s *serviceProvider) CacheUserRepository() repository.UserRepository {
	if s.cacheUserRepository == nil {
		s.cacheUserRepository = redisUserRepository.NewRepository(s.RedisClient())
	}

	return s.cacheUserRepository
}

// AccessRoleRepository access role repository
func (s *serviceProvider) AccessRoleRepository(ctx context.Context) repository.AccessibleRoleRepository {
	if s.accessibleRoleRepository == nil {
		s.accessibleRoleRepository = accessibleRoleRepository.NewRepository(s.DBClient(ctx))
	}

	return s.accessibleRoleRepository
}

// AuditLogService audit log service
func (s *serviceProvider) AuditLogService(ctx context.Context) service.AuditLogService {
	if s.auditLogService == nil {
		s.auditLogService = auditLogService.NewService(
			s.LogRepository(ctx),
		)
	}

	return s.auditLogService
}

// HashService hashing service
func (s *serviceProvider) HashService(ctx context.Context) service.HashService {
	if s.hashService == nil {
		s.hashService = hashService.NewService(
			s.HashingConfig().Salt(ctx),
		)
	}

	return s.hashService
}

// AuthService auth service
func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.UserRepository(ctx),
			s.AccessRoleRepository(ctx),
			s.ValidatorService(ctx),
			s.AuditLogService(ctx),
			s.HashService(ctx),
			s.CacheUserService(),
			s.TxManager(ctx),
			s.ProducerService(),
			s.AsyncRunner(),
			s.NewTokenGenerator(),
			s.NewAuthConfig(),
		)
	}

	return s.authService
}

// CacheUserService cache user service
func (s *serviceProvider) CacheUserService() service.CacheUserService {
	if s.cacheUserService == nil {
		s.cacheUserService = cacheUserService.NewService(
			s.RedisClient(),
			s.CacheUserRepository(),
		)
	}

	return s.cacheUserService
}

// ValidatorService validator service
func (s *serviceProvider) ValidatorService(ctx context.Context) service.ValidatorService {
	if s.validatorService == nil {
		s.validatorService = validator.NewValidator(s.UserRepository(ctx))
	}

	return s.validatorService
}

// ConsumerService kafka consumer service
func (s *serviceProvider) ConsumerService(ctx context.Context) service.ConsumerService {
	if s.consumerService == nil {
		s.consumerService = userSaverConsumer.NewService(
			s.AuthService(ctx),
			s.Consumer(),
		)
	}

	return s.consumerService
}

// Consumer kafka consumer
func (s *serviceProvider) Consumer() kafka.Consumer {
	if s.consumer == nil {
		s.consumer = kafkaConsumer.NewConsumer(
			s.ConsumerGroup(),
			s.ConsumerGroupHandler(),
		)
		closer.Add(s.consumer.Close)
	}

	return s.consumer
}

// ConsumerGroup kafka consumer group
func (s *serviceProvider) ConsumerGroup() sarama.ConsumerGroup {
	if s.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			s.KafkaConsumerConfig().Brokers(),
			s.KafkaConsumerConfig().GroupID(),
			s.KafkaConsumerConfig().Config(),
		)
		if err != nil {
			log.Fatalf("failed to create consumer group: %v", err)
		}

		s.consumerGroup = consumerGroup
	}

	return s.consumerGroup
}

// ConsumerGroupHandler consumer group handler
func (s *serviceProvider) ConsumerGroupHandler() *kafkaConsumer.GroupHandler {
	if s.consumerGroupHandler == nil {
		s.consumerGroupHandler = kafkaConsumer.NewGroupHandler()
	}

	return s.consumerGroupHandler
}

// SyncProducer kafka sync producer
func (s *serviceProvider) SyncProducer() kafka.SyncProducer {
	if s.syncProducer == nil {
		pr, err := producer.NewSyncProducer(
			s.KafkaProducerConfig(),
		)
		if err != nil {
			log.Fatalf("failed to create producer: %v", err)
		}
		s.syncProducer = pr
		closer.Add(s.syncProducer.Close)
	}

	return s.syncProducer
}

// ProducerService kafka producer service
func (s *serviceProvider) ProducerService() service.ProducerService {
	if s.producerService == nil {
		s.producerService = auditLogSender.NewService(
			s.SyncProducer(),
		)
	}

	return s.producerService
}

// AuthImpl auth api implementation
func (s *serviceProvider) AuthImpl(ctx context.Context) *auth.Implementation {
	if s.authImpl == nil {
		s.authImpl = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImpl
}
