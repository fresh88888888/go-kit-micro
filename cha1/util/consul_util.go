package util

import (
	"log"
	"strconv"

	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
)

var consulClient *consulapi.Client
var serviceId string
var serviceName string
var ServicePort int

func init() {
	config := consulapi.DefaultConfig()
	config.Address = "198.19.37.126:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	consulClient = client
	serviceId = "userService" + uuid.New().String()
}

func SetNameAndPort(name string, port int) {
	serviceName = name
	ServicePort = port
}

func RegService() {

	reg := consulapi.AgentServiceRegistration{}
	reg.ID = serviceId
	reg.Name = serviceName
	reg.Address = "198.19.37.126"
	reg.Port = ServicePort
	reg.Tags = []string{"primary"}

	check := consulapi.AgentServiceCheck{}
	check.HTTP = "http://198.19.37.126:" + strconv.Itoa(ServicePort) + "/health"
	check.Interval = "30s"

	reg.Check = &check

	err := consulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		log.Fatal(err)
	}
}

func UnRegService() {
	consulClient.Agent().ServiceDeregister(serviceId)
}
