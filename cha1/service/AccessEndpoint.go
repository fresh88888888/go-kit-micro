package service

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
)

const secKey = "123abc"

type UserClaim struct {
	UserName string `json:"username"`
	jwt.StandardClaims
}

type IAccessService interface {
	GetToken(userName string, upass string) (string, error)
}

type AccessService struct{}

func (a *AccessService) GetToken(userName string, upass string) (string, error) {
	if userName == "zcj" && upass == "123456" {
		userInfo := &UserClaim{UserName: userName}
		userInfo.ExpiresAt = time.Now().Add(time.Second * 20).Unix()
		tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, userInfo)
		token, err := tokenObj.SignedString([]byte(secKey))
		return token, err
	}

	return "", fmt.Errorf("error username and password")
}

type AccessRequest struct {
	UserName string
	Password string
	Method   string
}

type AccessResponse struct {
	Status string
	Token  string
}

func AccessEndpoint(accessService IAccessService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		r := request.(AccessRequest)
		result := AccessResponse{Status: "ok"}
		if r.Method == "POST" {
			token, err := accessService.GetToken(r.UserName, r.Password)
			if err != nil {
				result.Status = "error:" + err.Error()
			} else {
				result.Token = token
			}
		}

		return result, nil
	}
}
