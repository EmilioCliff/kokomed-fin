package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/handlers"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/payments"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/redis"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/services"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/workers"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

func main() {
	// Setup signal handlers.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	config, err := pkg.LoadConfig("../../.envs/.local", "config", "yaml")
	if err != nil {
		panic(err)
	}
	log.Println("config loaded successfuly")

	maker, err := pkg.NewJWTMaker(config.RSA_PRIVATE_KEY, config.RSA_PUBLIC_KEY, config)
	if err != nil {
		panic(err)
	}
	log.Println("JWT maker loaded successfuly")

	store := mysql.NewStore(config)

	err = store.OpenDB()
	if err != nil {
		panic(err)
	}
	log.Println("db opened successfuly")

	cache := redis.NewCacheClient(config.REDIS_ADDRESS, config.REDIS_PASSWORD, 1)

	repo := mysql.NewMySQLRepo(store)
	paymentService := payments.NewPaymentService(repo, store)

	sender := pkg.NewGmailSender(config.EMAIL_SENDER_NAME, config.EMAIL_SENDER_ADDRESS, config.EMAIL_SENDER_PASSWORD)

	redisConfig := services.RedisConfig{
		Address: config.REDIS_ADDRESS,
		Password: config.REDIS_PASSWORD,
		DB: 0,
	}

	worker := workers.NewWorkerService(redisConfig, sender, repo, *maker)

	err = worker.StartProcessor()
	if err != nil {
		panic(err)
	}

	log.Println("Started worker successfuly")


	server := handlers.NewServer(config, *maker, repo, paymentService, worker, cache)

	log.Println("starting server at port ", config.HTTP_PORT)
	if err := server.Start(); err != nil {
		panic(err)
	}

	token, _ := maker.CreateToken("emiliocliff@gmail.com", 1, 1, "ADMIN", 10*time.Hour)

	log.Println(token)

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
