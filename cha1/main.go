package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	httptransport "github.com/go-kit/kit/transport/http"
	kitlog "github.com/go-kit/log"
	"github.com/gorilla/mux"
	"golang.org/x/time/rate"
	"umbrella.github.com/go-kit-example/cha1/service"
	"umbrella.github.com/go-kit-example/cha1/util"
)

func main() {
	name := flag.String("n", "", "service name")
	port := flag.Int("p", 0, "service port")
	flag.Parse()
	if *name == "" {
		log.Fatal("Please input service name: ")
	}
	if *port == 0 {
		log.Fatal("Please input service port: ")
	}
	util.SetNameAndPort(*name, *port)

	var logger kitlog.Logger
	{
		logger = kitlog.NewLogfmtLogger(os.Stdout)
		logger = kitlog.WithPrefix(logger, "mykit", "v1")
		logger = kitlog.With(logger, "time", kitlog.DefaultTimestamp)
		logger = kitlog.With(logger, "invoker", kitlog.DefaultCaller)
	}

	user := service.UserService{}
	rate_limiter := rate.NewLimiter(1, 5)
	endp := service.RateLimit(rate_limiter)(service.UserServiceLog(logger)(service.CheckToken()(service.GenUserEndpoint(user))))
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(service.MyErrorEncoder),
	}

	serverhandler := httptransport.NewServer(endp, service.DecodeUserRequest, service.EncodeUserResponse, options...)

	access := &service.AccessService{}
	accessEndpoint := service.AccessEndpoint(access)
	accessHandler := httptransport.NewServer(accessEndpoint, service.DecodeAccessRequest, service.EncodeAcessResponse, options...)

	r := mux.NewRouter()
	r.Methods("POST").Path("/access-token").Handler(accessHandler)
	r.Methods("GET", "DELETE").Path(`/user/{uid:\d+}`).Handler(serverhandler)
	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/health", checkHealth).Methods("GET")
	// run http service
	error_channel := make(chan error)
	go func() {
		util.RegService()
		if err := http.ListenAndServe(":"+strconv.Itoa(*port), r); err != nil {
			log.Println(err.Error())
			error_channel <- err
		}
	}()

	go func() {
		sig_channel := make(chan os.Signal, 1)
		signal.Notify(sig_channel, syscall.SIGINT, syscall.SIGTERM)
		error_channel <- fmt.Errorf("%s", <-sig_channel)
	}()

	s_err := <-error_channel
	util.UnRegService()
	log.Println(s_err.Error())

}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Home page.")
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{"status": "ok"}`))
}
