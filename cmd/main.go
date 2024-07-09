package main

import (
	"context"
	"github.com/Lidne/praktika_MAI/config"
	_ "github.com/Lidne/praktika_MAI/docs"
	"github.com/Lidne/praktika_MAI/internal/server"
	"github.com/Lidne/praktika_MAI/pkg/jaeger"
	"github.com/Lidne/praktika_MAI/pkg/kafka"
	"github.com/Lidne/praktika_MAI/pkg/logger"
	"github.com/Lidne/praktika_MAI/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opentracing/opentracing-go"
	"log"
)

type App struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewApp(pool *pgxpool.Pool) *App {
	return &App{pool: pool, ctx: context.Background()}
}

// @title           Stats microservice
// @version         1.0
// @description     Statistics microservice
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:5007
// @BasePath  /api/

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
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

	dbpool, err := postgres.NewClient(ctx, cfg)
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

	s := server.NewServer(appLogger, cfg, tracer, dbpool)
	appLogger.Fatal(s.Run())
}
