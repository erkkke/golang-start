package main

import (
	"context"
	"github.com/erkkke/golang-start/project/internal/cache/redis_cache"
	"github.com/erkkke/golang-start/project/internal/http"
	"github.com/erkkke/golang-start/project/internal/store/postgres"
	"log"
)

const (
	port = ":8081"

	cacheDB = 1
	cacheExpTime = 1800
	cachePort = "localhost:6379"
)

func main() {
	urlDB := "postgres://postgres:postgres@localhost:5432/golang_project"
	store := postgres.NewDB()
	if err := store.Connect(urlDB); err != nil {
		panic(err)
	}
	defer store.Close()


	cache := redis_cache.NewRedisCache(cachePort, cacheDB, cacheExpTime)
	//defer cache.Close()

	srv := http.NewServer(context.Background(),
		http.WithAddress(port),
		http.WithStore(store),
		http.WithCache(cache),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
