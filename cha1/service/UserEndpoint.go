package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"golang.org/x/time/rate"
	"umbrella.github.com/go-kit-example/cha1/util"
)

type UserRequest struct {
	UserId int `json:"uid"`
	Method string
	Token  string
}

type UserResponse struct {
	Result string `json:"data"`
}

func CheckToken() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			uc := UserClaim{}
			get_token, err := jwt.ParseWithClaims(r.Token, &uc, func(t *jwt.Token) (interface{}, error) {
				return []byte(secKey), nil
			})

			if err != nil {
				return nil, util.NewServerError(430, err.Error())
			}

			if get_token != nil && get_token.Valid {
				c := context.WithValue(ctx, "login_user", get_token.Claims.(*UserClaim).UserName)
				return next(c, request)
			} else {
				return nil, util.NewServerError(430, "token error")
			}
		}
	}
}

func UserServiceLog(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			r := request.(UserRequest)
			logger.Log("method", r.Method, "event", "get_user", "userid", r.UserId)
			return next(ctx, request)
		}
	}
}

func RateLimit(limit *rate.Limiter) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			if !limit.Allow() {
				return nil, util.NewServerError(429, "too many reques")
			}
			return next(ctx, request)
		}
	}
}
func GenUserEndpoint(userService IUserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(UserRequest)
		fmt.Println("name of current login user: ", ctx.Value("login_user"))
		result := "nothing"
		if r.Method == "GET" {
			result = userService.GetName(r.UserId) + strconv.Itoa(util.ServicePort)
		} else if r.Method == "DELETE" {
			err := userService.DelUser(r.UserId)
			if err != nil {
				result = err.Error()
			} else {
				result = fmt.Sprintf("the user delete success, for userid:%d", r.UserId)
			}
		}

		return UserResponse{Result: result}, nil
	}
}
