package main

import (
	"context"
	"github.com/erkkke/golang-start/hw6/internal/http"
	"github.com/erkkke/golang-start/hw6/internal/store/postgres"
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

	srv := http.NewServer(context.Background(), port, store)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
