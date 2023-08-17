package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func GetUserInfoReq(_ context.Context, r *http.Request, i interface{}) error {
	user_req := i.(UserRequest)
	r.URL.Path = "/user/" + strconv.Itoa(user_req.UserId)
	return nil
}

func GetUserInfoRep(_ context.Context, r *http.Response) (response interface{}, err error) {
	if r.StatusCode >= 400 {
		return nil, errors.New("no data")
	}

	var user_reqponse UserResponse
	err = json.NewDecoder(r.Body).Decode(&user_reqponse)
	if err != nil {
		return nil, err
	}

	return user_reqponse, nil
}
