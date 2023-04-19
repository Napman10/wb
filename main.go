package main

import (
	"context"
	"log"
	"os"

	"wb/database"
	"wb/service"
	"wb/transport"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("POSTGRES_DSN")

	db, err := database.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
		return
	}

	port := os.Getenv("HTTP_PORT")

	_, err = transport.New(service.New(db), port)
	if err != nil {
		log.Fatal(err)
		return
	}
}
