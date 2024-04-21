package user

import (
	"context"
	"errors"
	"log"

	"github.com/germansanz93/go-web-fundamentals/internal/user/domain"
	"golang.org/x/exp/slices"
)

type DB struct {
	Users     []domain.User
	MaxUserId uint64
}

type (
	Repository interface {
		Create(ctx context.Context, user *domain.User) error
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
		Update(ctx context.Context, id uint64, firstName, lastName, email *string) error
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

func (r *repo) Get(ctx context.Context, id uint64) (*domain.User, error) {
	index := slices.IndexFunc(r.db.Users, func(v domain.User) bool {
		return v.ID == id
	})

	if index < 0 {
		return nil, errors.New("user not found")
	}

	return &r.db.Users[index], nil
}

func (r *repo) Update(ctx context.Context, id uint64, firstName, lastName, email *string) error {
	user, err := r.Get(ctx, id)
	if err != nil {
		return err
	}
	if firstName != nil {
		user.FirstName = *firstName
	}
	if lastName != nil {
		user.LastName = *lastName
	}
	if email != nil {
		user.Email = *email
	}
	return nil
}
