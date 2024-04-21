package user

import (
	"context"
	"log"

	"github.com/germansanz93/go-web-fundamentals/internal/user/domain"
)

type (
	Service interface {
		Create(ctx context.Context, firsName, lastName, email string) (*domain.User, error)
		GetAll(ctx context.Context) ([]domain.User, error)
		Get(ctx context.Context, id uint64) (*domain.User, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

func NewService(l *log.Logger, repo Repository) Service {
	return &service{
		log:  l,
		repo: repo,
	}
}

func (s service) Create(ctx context.Context, firsName, lastName, email string) (*domain.User, error) {
	user := &domain.User{
		FirstName: firsName,
		LastName:  lastName,
		Email:     email,
	}
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	s.log.Println("service: user create")
	return user, nil //Tener en cuenta que al estar laburando con puntero, al el repository modificarlo, estaremos devolviendo el modificado
}

func (s service) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}
	s.log.Println("service: users get all")
	return users, nil
}

func (s service) Get(ctx context.Context, id uint64) (*domain.User, error) {
	return s.repo.Get(ctx, id)
}
