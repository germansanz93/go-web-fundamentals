package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/germansanz93/go-web-fundamentals/internal/user"
	"github.com/germansanz93/go-web-fundamentals/pkg/bootstrap"
	"github.com/germansanz93/go-web-fundamentals/pkg/handler"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()

	server := http.NewServeMux()

	db, err := bootstrap.NewDB()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	logger := bootstrap.NewLogger()
	repo := user.NewRepo(db, logger)

	service := user.NewService(logger, repo)

	ctx := context.Background()

	handler.NewUserHttpServer(ctx, server, user.MakeEndpoints(ctx, service))

	port := os.Getenv("PORT")
	log.Println("Server started at port: ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}
