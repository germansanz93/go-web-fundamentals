package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/germansanz93/go-web-fundamentals/internal/user"
	"github.com/germansanz93/go-web-fundamentals/pkg/transport"
)

func NewUserHttpServer(ctx context.Context, router *http.ServeMux, endpoints user.Endpoints) {
	router.HandleFunc("/users", UserServer(ctx, endpoints))

}

func UserServer(ctx context.Context, endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		tran := transport.New(w, r, ctx)

		switch r.Method {
		case http.MethodGet:
			tran.Server(
				transport.Endpoint(endpoints.GetAll),
				decodeGetAllUser,
				encodeResponse,
				encodeError,
			)
			return
		case http.MethodPost:
			tran.Server(
				transport.Endpoint(endpoints.Create),
				decodeCreateUser,
				encodeResponse,
				encodeError,
			)
			return
		default:
			InvalidMethod(w)
		}
	}
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	status := http.StatusInternalServerError
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "%s"}`, status, err)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	status := http.StatusOK
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "data": %s}`, status, data)
	return nil
}

func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("invalid request format: '%v'", err.Error())
	}
	return req, nil
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exist"}`, status)

}
