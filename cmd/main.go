package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/germansanz93/go-web-fundamentals/internal/user"
	"github.com/germansanz93/go-web-fundamentals/internal/user/domain"
)

func main() {
	server := http.NewServeMux()

	db := user.DB{
		Users: []domain.User{
			{
				ID:        1,
				FirstName: "German",
				LastName:  "Sanz",
				Email:     "germansanz@gmail.com",
			},
			{
				ID:        2,
				FirstName: "Jhon",
				LastName:  "Doe",
				Email:     "johndoe@gmail.com",
			},
			{
				ID:        3,
				FirstName: "Mister",
				LastName:  "Jagger",
				Email:     "jaggermister@gmail.com",
			},
		},
		MaxUserId: 3,
	}

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	repo := user.NewRepo(db, logger)

	service := user.NewService(logger, repo)

	ctx := context.Background()

	server.HandleFunc("/users", user.MakeEndpoints(ctx, service))
	log.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
