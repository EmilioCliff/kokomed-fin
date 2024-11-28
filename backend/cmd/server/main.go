package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EmilioCliff/kokomed-fin/backend/internal/handlers"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/mysql"
	"github.com/EmilioCliff/kokomed-fin/backend/internal/payments"
	"github.com/EmilioCliff/kokomed-fin/backend/pkg"
)

func main() {
	// Setup signal handlers.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	config, err := pkg.LoadConfig("../../.envs/.local")
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

	repo := mysql.NewMySQLRepo(store)
	paymentService := payments.NewPaymentService(repo, store)

	server := handlers.NewServer(config, *maker, repo, paymentService)

	log.Println("starting server at port ", config.HTTP_PORT)
	if err := server.Start(); err != nil {
		panic(err)
	}

	<-quit

	if err = store.CloseDB(); err != nil {
		log.Fatal(err)
	}

	if err = server.Stop(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server shutdown ...")
}
