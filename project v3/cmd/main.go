package main

import (
	"context"
	"github.com/erkkke/golang-start/project/internal/http"
	"github.com/erkkke/golang-start/project/internal/message_broker/kafka"
	"github.com/erkkke/golang-start/project/internal/store/postgres"
	"github.com/erkkke/golang-start/project/pkg/auth"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = ":8081"
const key = "secret"


func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	dbURL := "postgres://postgres:postgres@localhost:5432/golang_project"
	store := postgres.NewDB()
	if err := store.Connect(dbURL); err != nil {
		panic(err)
	}
	defer store.Close()

	cache, err := lru.New2Q(10)
	if err != nil {
		panic(err)
	}

	manager, err := auth.NewManager(key)
	if err != nil {
		panic(err)
	}

	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer3")
	if err = broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	srv := http.NewServer(context.Background(),
		http.WithAddress(port),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithBroker(broker),
		http.WithTokenManager(manager),
	)
	if err = srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}
