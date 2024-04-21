package bootstrap

import (
	"log"
	"os"

	"github.com/germansanz93/go-web-fundamentals/internal/user"
	"github.com/germansanz93/go-web-fundamentals/internal/user/domain"
)

func NewLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func NewDB() user.DB {
	return user.DB{
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
}
