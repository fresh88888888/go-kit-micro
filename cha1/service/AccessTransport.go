package service

import (
	"context"
	"encoding/json"
	"net/http"
)

type Result struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func DecodeAccessRequest(c context.Context, req *http.Request) (interface{}, error) {
	var result Result
	decoder := json.NewDecoder(req.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return AccessRequest{
		UserName: result.UserName,
		Password: result.Password,
		Method:   req.Method,
	}, nil
}

func EncodeAcessResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}
