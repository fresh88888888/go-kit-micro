package util

import (
	"context"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	consulapi "github.com/hashicorp/consul/api"
	"umbrella.github.com/go-kit-example/cha2/service"
)

func GetUser() (string, error) {
	config := consulapi.DefaultConfig()
	config.Address = "198.19.37.126:8500"
	api_client, err := consulapi.NewClient(config)
	if err != nil {
		return "", err
	}

	client := consul.NewClient(api_client)

	var logger log.Logger = log.NewLogfmtLogger(os.Stdout)
	tags := []string{"primary"}
	instance := consul.NewInstancer(client, logger, "userservice", tags, true)

	factory := func(instance string) (endpoint.Endpoint, io.Closer, error) {
		url, _ := url.Parse("http://" + instance)
		return httptransport.NewClient("GET", url, service.GetUserInfoReq, service.GetUserInfoRep).Endpoint(), nil, nil
	}
	endpoint := sd.NewEndpointer(instance, factory, logger)
	// blancer := lb.NewRoundRobin(endpoint)
	blancer := lb.NewRandom(endpoint, time.Now().UnixNano())

	getUserInfo, err := blancer.Endpoint()
	if err != nil {
		return "", err
	}
	response, err := getUserInfo(context.Background(), service.UserRequest{UserId: 102})
	if err != nil {
		return "", err
	}
	userInfo := response.(service.UserResponse)

	time.Sleep(time.Second * 1)

	return userInfo.Result, nil
}
