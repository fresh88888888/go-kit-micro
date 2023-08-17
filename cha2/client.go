package main

import (
	"fmt"
	"time"

	"github.com/afex/hystrix-go/hystrix"
	"umbrella.github.com/go-kit-example/cha2/util"
)

func main() {
	// target, _ := url.Parse("http://198.19.37.126:8081")
	// client := httptransport.NewClient("GET", target, service.GetUserInfoReq, service.GetUserInfoRep)
	// getUserInfo := client.Endpoint()

	// response, err := getUserInfo(context.Background(), service.UserRequest{UserId: 101})
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }

	// fmt.Println(response.(service.UserResponse).Result)

	// fuse config items
	fause_config := hystrix.CommandConfig{
		Timeout:                2000,
		MaxConcurrentRequests:  5,
		RequestVolumeThreshold: 3,
		ErrorPercentThreshold:  20,
		SleepWindow:            int(time.Second * 2),
	}

	hystrix.ConfigureCommand("getUser", fause_config)
	for {
		err := hystrix.Do("getUser", func() error {
			res, err := util.GetUser()
			fmt.Println(res)
			return err
		}, func(e error) error {
			fmt.Println("user service degraded")
			return e
		})

		if err != nil {
			fmt.Println(err.Error())
			//log.Fatal(err)
		}
	}
}
