package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/handlers"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/otelImp"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/payments"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/redis"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/reports"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/workers"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// setup signal handlers.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// load config
	configPath := "/home/emilio-cliff/kokomed-fin/backend/.envs/.local"
	config, err := pkg.LoadConfig(configPath, "config", "yaml")
	if err != nil {
		panic(err)
	}

	// create a grpc client for sending telementry data to otel collector
	conn, err := grpc.NewClient("localhost:4317", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	log.Println("Connection state: ", conn.GetState())

	// create an instance of openTelemetry and initialize the providers
	openTelemetry, err := otelImp.NewOpenTelemetry()
	if err != nil {
		panic(err)
	}

	shutdownTracer, err := openTelemetry.InitializeTracerProvider(context.Background(), conn)
	if err != nil {
		panic(err)
	}

	shutdownMeter, err := openTelemetry.InitializeMeterProvider(context.Background(), conn)
	if err != nil {
		panic(err)
	}

	shutdownLogger, err := openTelemetry.InitializeLoggerProvider(context.Background(), conn)
	if err != nil {
		panic(err)
	}

	// defer the shutdown of the otel providers
	defer func() {
		if err := shutdownTracer(context.Background()); err != nil {
			log.Fatal(err)
		}
		if err := shutdownMeter(context.Background()); err != nil {
			log.Fatal(err)
		}
		if err := shutdownLogger(context.Background()); err != nil {
			log.Fatal(err)
		}
	}()

	// create a new store and connect to db
	store := mysql.NewStore(config)
	err = store.OpenDB()
	if err != nil {
		panic(err)
	}

	// initialize the our dependancies and start the server
	maker, err := pkg.NewJWTMaker(config.RSA_PRIVATE_KEY, config.RSA_PUBLIC_KEY, config)
	if err != nil {
		panic(err)
	}

	sender := pkg.NewGmailSender(config.EMAIL_SENDER_NAME, config.EMAIL_SENDER_ADDRESS, config.EMAIL_SENDER_PASSWORD)

	redisConfig := services.RedisConfig{
		Address: config.REDIS_ADDRESS,
		Password: config.REDIS_PASSWORD,
		DB: 0,
	}

	repo := mysql.NewMySQLRepo(store)

	worker := workers.NewWorkerService(redisConfig, sender, repo, *maker)
	err = worker.StartProcessor()
	if err != nil {
		panic(err)
	}

	report := reports.NewReportService(repo)
	paymentService := payments.NewPaymentService(repo, store)
	cache := redis.NewCacheClient(config.REDIS_ADDRESS, config.REDIS_PASSWORD, 1)

	server := handlers.NewServer(config, *maker, repo, paymentService, worker, cache, report)

	log.Println("starting server at port: ", config.HTTP_PORT)
	if err := server.Start(); err != nil {
		panic(err)
	}

	// only for testing
	token, _ := maker.CreateToken("emiliocliff@gmail.com", 1, 1, "ADMIN", 10*time.Hour)
	log.Println(token)

	// get cpu metrics
	unregisterCPUMetrics, err := pkg.GetCPUMetrics()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := unregisterCPUMetrics.Unregister(); err != nil {
			log.Fatal(err)
		}
	}()

	<-quit

	if err = store.CloseDB(); err != nil {
		log.Fatal(err)
	}

	if err = server.Stop(); err != nil {
		log.Fatal(err)
	}

	worker.StopProcessor()

	log.Println("Server shutdown ...")
}
