package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/germansanz93/go-fundamentals-response/response"
	"github.com/germansanz93/go-web-fundamentals/internal/user"
	"github.com/germansanz93/go-web-fundamentals/pkg/transport"
	"github.com/gin-gonic/gin"
)

func NewUserHttpServer(endpoints user.Endpoints) http.Handler {
	r := gin.Default()

	r.POST("/users", transport.GinServer(
		transport.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
		encodeError,
	))
	r.GET("/users", transport.GinServer(
		transport.Endpoint(endpoints.GetAll),
		decodeGetAllUser,
		encodeResponse,
		encodeError,
	))
	r.GET("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		encodeError,
	))
	r.PATCH("/users/:id", transport.GinServer(
		transport.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		encodeError,
	))
	return r
}

func encodeError(c *gin.Context, err error) {
	c.Header("Content-Type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	c.JSON(resp.StatusCode(), resp)
}

func encodeResponse(c *gin.Context, resp interface{}) {
	r := resp.(response.Response)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(r.StatusCode(), resp)
}

func decodeCreateUser(c *gin.Context) (interface{}, error) {
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	var req user.CreateReq
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

func decodeUpdateUser(c *gin.Context) (interface{}, error) {
	var req user.UpdateReq
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	if err := tokenVerify(c.Request.Header.Get("Authorization")); err != nil {
		return nil, response.Unauthorized(err.Error())
	}
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
	if err != nil {
		return nil, response.BadRequest(err.Error())
	}

	req.Id = id

	return req, nil
}

func decodeGetAllUser(c *gin.Context) (interface{}, error) {
	return nil, nil
}

func decodeGetUser(c *gin.Context) (interface{}, error) {
	id, err := strconv.ParseUint(c.Params.ByName("id"), 10, 64)
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
