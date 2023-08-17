package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"umbrella.github.com/go-kit-example/cha1/util"
)

func DecodeUserRequest(c context.Context, req *http.Request) (interface{}, error) {
	vars := mux.Vars(req)
	if uid, ok := vars["uid"]; ok {
		uid, _ := strconv.Atoi(uid)
		return UserRequest{
			UserId: uid,
			Method: req.Method,
			Token:  req.URL.Query().Get("token"),
		}, nil
	}
	return nil, errors.New("parameter error")
}

func EncodeUserResponse(c context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

func MyErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	contentType, body := "text/plain; charet=utf-8", []byte(err.Error())
	w.Header().Set("Content-type", contentType)
	if myerr, ok := err.(*util.ServerError); ok {
		w.WriteHeader(myerr.Code)
		w.Write([]byte(myerr.Msg))
		return
	} else {
		w.WriteHeader(500)
		w.Write(body)
	}
}
