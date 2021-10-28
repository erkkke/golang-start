package main

import (
	"context"
	"github.com/erkkke/golang-start/hw6/internal/http"
	"github.com/erkkke/golang-start/hw6/internal/store/inmemory"
	"log"
)

const port = ":8081"

func main() {
	store := inmemory.NewDB()

	srv := http.NewServer(context.Background(), port, store)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}
