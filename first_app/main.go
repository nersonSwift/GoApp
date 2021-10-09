package main

import (
	"MyProjects/go_app/models/appConfig"
	"MyProjects/go_app/models/dbModels"
	"MyProjects/go_app/servises/rmqServise"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

var config = appConfig.Get()
var rmq = rmqServise.Connect(config.RMQConnectLink)

var rmqConnect *amqp.Connection

func main() {
	fmt.Println("Start")

	createCreateHandler1()
	//createCreateHandler()
	//createHelloHandler()
	c1 := make(chan string)
	<-c1
}

func createHelloHandler() {
	handler := http.NewServeMux()
	handler.HandleFunc("/go/hello", helloHandler)

	s := http.Server{
		Addr:           ":70",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())
}

func createCreateHandler() {
	handler := http.NewServeMux()
	handler.HandleFunc("/go/create", createHandler)
	handler.HandleFunc("/go/hello", helloHandler)

	s := http.Server{
		Addr:           ":70",
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	handler.HandleFunc("/go/hello1", helloHandler)

	log.Fatal(s.ListenAndServe())

	handler.HandleFunc("/go/hello2", helloHandler)

	log.Fatal(s.ListenAndServe())
}

func createCreateHandler1() {
	http.HandleFunc("/go/create/", createHandler)
	http.HandleFunc("/go/hello", helloHandler)
	go func() {
		http.ListenAndServe(":"+config.Port, nil)
	}()
}

// Resp is json model
type Resp struct {
	Message string
	Error   string
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("1 zapros")
	//createRmqConnect()
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	//createRmqConnect()
	mess := dbModels.User{
		ID:      1,
		Name:    "Test",
		Country: "ru",
	}
	rmq.SendMassage(mess)

	fmt.Println("1 zapros")
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()

	for key, value := range r.Form {
		if key == "key" {
			http.HandleFunc("/go/"+value[0], helloHandler)
		}
		fmt.Println(key, value)
	}
	fsdf := PostP(config.SecondLink + "/go/testAPI")
	fmt.Println(config.SecondLink)
	resp := Resp{
		Message: fsdf,
	}
	respJ, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respJ)
}

// PostP is ...
func PostP(path string) string {
	data := url.Values{}

	req, err := http.NewRequest("GET", path, strings.NewReader(data.Encode()))
	if err != nil {
		var errB []byte
		return string(errB)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		var errB []byte
		return string(errB)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
