package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/germansanz93/go-fundamentals-response/response"
	"github.com/germansanz93/go-web-fundamentals/internal/user"
	"github.com/germansanz93/go-web-fundamentals/pkg/transport"
)

func NewUserHttpServer(ctx context.Context, router *http.ServeMux, endpoints user.Endpoints) {
	router.HandleFunc("/users/", UserServer(ctx, endpoints))

}

func UserServer(ctx context.Context, endpoints user.Endpoints) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		log.Println(r.Method, url)

		path, pathSize := transport.Clean(url)

		params := make(map[string]string)
		if pathSize == 4 && path[2] != "" {
			params["userId"] = path[2]
		}

		params["token"] = r.Header.Get("Authorization")

		tran := transport.New(w, r, context.WithValue(ctx, "params", params))

		var end user.Controller
		var deco func(ctx context.Context, r *http.Request) (interface{}, error)

		switch r.Method {
		case http.MethodGet:
			switch pathSize {
			case 3:
				end = endpoints.GetAll
				deco = decodeGetAllUser
			case 4:
				end = endpoints.Get
				deco = decodeGetUser
			}

		case http.MethodPost:
			switch pathSize {
			case 3:
				end = endpoints.Create
				deco = decodeCreateUser
			}

		case http.MethodPatch:
			switch pathSize {
			case 4:
				end = endpoints.Update
				deco = decodeUpdateUser
			}

		default:
			InvalidMethod(w)
		}

		if end != nil && deco != nil {
			tran.Server(
				transport.Endpoint(end),
				deco,
				encodeResponse,
				encodeError,
			)
		} else {
			InvalidMethod(w)
		}
	}
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(resp)
}

func decodeCreateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	params := ctx.Value("params").(map[string]string)
	if err := tokenVerify(params["token"]); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

func decodeUpdateUser(ctx context.Context, r *http.Request) (interface{}, error) {
	var req user.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	params := ctx.Value("params").(map[string]string)
	if err := tokenVerify(params["token"]); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	id, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	req.Id = id

	return req, nil
}

func decodeGetAllUser(ctx context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func decodeGetUser(ctx context.Context, r *http.Request) (interface{}, error) {
	params := ctx.Value("params").(map[string]string)
	id, err := strconv.ParseUint(params["userId"], 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}
	return user.GetReq{Id: id}, nil
}

func tokenVerify(token string) error {
	if os.Getenv("TOKEN") != token {
		return errors.New("invalid token")
	}
	return nil
}

func InvalidMethod(w http.ResponseWriter) {
	status := http.StatusNotFound
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"status": %d, "message": "method doesn't exist"}`, status)

}
