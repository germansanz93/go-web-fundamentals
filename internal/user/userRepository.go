package user

import (
	"context"
	"log"

	"github.com/germansanz93/go-web-fundamentals/internal/user/domain"
)

type DB struct {
	Users     []domain.User
	MaxUserId uint64
}

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
	}

	repo struct {
		db  DB
		log *log.Logger
	}
)

func NewRepo(db DB, l *log.Logger) Repository {
	return &repo{
		db:  db,
		log: l,
	}
}

func (r *repo) Create(ctx context.Context, user *domain.User) error {
	r.db.MaxUserId++
	user.ID = r.db.MaxUserId
	r.db.Users = append(r.db.Users, *user)
	r.log.Println("repository: user create")
	return nil
}

func (r *repo) GetAll(ctx context.Context) ([]domain.User, error) {
	r.log.Println("repository: user get all")
	return r.db.Users, nil
}
