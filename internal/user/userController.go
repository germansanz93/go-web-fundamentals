package user

import (
	"context"
)

type (
	Controller func(ctx context.Context, request interface{}) (interface{}, error)

	Endpoints struct {
		Create Controller
		GetAll Controller
		Get    Controller
		Update Controller
	}

	GetReq struct {
		Id uint64
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	UpdateReq struct {
		Id        uint64
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
	}
)

func MakeEndpoints(ctx context.Context, s Service) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		GetAll: makeGetAllEndpoint(s),
		Get:    makeGetEndpoint(s),
		Update: makeUpdateEndpoint(s),
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

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetReq)
		return s.Get(ctx, req.Id)
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		userRequest := request.(CreateReq) //Casteo la interfaz a User

		if userRequest.FirstName == "" {
			return nil, ErrFirstNameRequired
		}

		if userRequest.LastName == "" {
			return nil, ErrLastNameRequired
		}

		if userRequest.Email == "" {
			return nil, ErrEmailRequired
		}
		user, err := s.Create(ctx, userRequest.FirstName, userRequest.LastName, userRequest.Email)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateReq)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, ErrFirstNameRequired
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, ErrLastNameRequired
		}

		if req.Email != nil && *req.Email == "" {
			return nil, ErrEmailRequired
		}

		if err := s.Update(ctx, req.Id, req.FirstName, req.LastName, req.Email); err != nil {
			return nil, err
		}
		return nil, nil
	}
}
