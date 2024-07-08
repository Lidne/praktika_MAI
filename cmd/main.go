package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"

	"github.com/opentracing/opentracing-go"

	"github.com/AleksK1NG/products-microservice/config"
	"github.com/AleksK1NG/products-microservice/internal/server"
	"github.com/AleksK1NG/products-microservice/pkg/jaeger"
	"github.com/AleksK1NG/products-microservice/pkg/kafka"
	"github.com/AleksK1NG/products-microservice/pkg/logger"
	"github.com/AleksK1NG/products-microservice/pkg/redis"
)

// @title Products microservice
// @version 1.0
// @description Products REST API
// @termsOfService http://swagger.io/terms/

// @contact.name Alexander Bryksin
// @contact.url https://github.com/AleksK1NG
// @contact.email alexander.bryksin@yandex.ru

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:5007
// @BasePath /api/v1

func main() {
	log.Println("Starting products microservice")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting user server")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Development,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	dbpool, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if pingErr := dbpool.Ping(ctx); err != nil || pingErr != nil {
		appLogger.Fatal("cannot connect mongodb", err)
	}
	defer dbpool.Close()

	appLogger.Info("PostgreSQL connected")

	conn, err := kafka.NewKafkaConn(cfg)
	if err != nil {
		appLogger.Fatal("NewKafkaConn", err)
	}
	defer conn.Close()
	brokers, err := conn.Brokers()
	if err != nil {
		appLogger.Fatal("conn.Brokers", err)
	}
	appLogger.Infof("Kafka connected: %v", brokers)

	redisClient := redis.NewRedisClient(cfg)
	appLogger.Info("Redis connected")

	s := server.NewServer(appLogger, cfg, tracer, dbpool, redisClient)
	appLogger.Fatal(s.Run())
}
