package main

import (
	"context"
	"github.com/erkkke/golang-start/project/internal/http"
	"github.com/erkkke/golang-start/project/internal/store/postgres"
	lru "github.com/hashicorp/golang-lru"
	"log"
)

const port = ":8081"

func main() {
	urlExample := "postgres://postgres:postgres@localhost:5432/golang_project"
	store := postgres.NewDB()
	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	cache, err := lru.New2Q(6)
	if err != nil {
		panic(err)
	}

	srv := http.NewServer(context.Background(), port, store, cache)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
