package main

import (
	"context"
	"log"
	"net/http"

	"github.com/germansanz93/go-web-fundamentals/internal/user"
	"github.com/germansanz93/go-web-fundamentals/pkg/bootstrap"
	"github.com/germansanz93/go-web-fundamentals/pkg/handler"
)

func main() {
	server := http.NewServeMux()

	db := bootstrap.NewDB()

	logger := bootstrap.NewLogger()
	repo := user.NewRepo(db, logger)

	service := user.NewService(logger, repo)

	ctx := context.Background()

	handler.NewUserHttpServer(ctx, server, user.MakeEndpoints(ctx, service))

	log.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
