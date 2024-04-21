package user

import (
	"context"
	"errors"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
)

func MakeEndpoints(ctx context.Context, s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
	}
}

func makeGetAllEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		users, err := s.GetAll(ctx)
		if err != nil {
			return nil, err
		}
		return users, nil
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userRequest := request.(CreateReq) //Casteo la interfaz a User

		if userRequest.FirstName == "" {
			return nil, errors.New("first name is required")
		}

		if userRequest.LastName == "" {
			return nil, errors.New("last name is required")
		}

		if userRequest.Email == "" {
			return nil, errors.New("email is required")
		}
		user, err := s.Create(ctx, userRequest.FirstName, userRequest.LastName, userRequest.Email)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}
