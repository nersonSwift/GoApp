package main

import (
	"MyProjects/go_app/appConfig"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

type Mess struct {
	Message string
}

var config = appConfig.Get()

func main() {
	fmt.Println("Start")
	createHelloHandler()

	c1 := make(chan string)
	<-c1
}

var rmqConnect *amqp.Connection

func createRmqConnect() {
	rmqConnectL, err := amqp.Dial(config.RMQConnectLink)
	if err != nil {
		fmt.Println(err)
		return
	}
	rmqConnect = rmqConnectL

	channel, err := rmqConnect.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	queue, err := channel.QueueDeclare("add", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = channel.Qos(1, 0, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	messageChannel, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	stopChan := make(chan bool)
	go func() {
		for d := range messageChannel {
			addTask := Mess{}
			err := json.Unmarshal(d.Body, &addTask)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Mess: " + addTask.Message)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
		fmt.Println("End")
	}()
	<-stopChan
}

func createHelloHandler() {
	http.HandleFunc("/go/testAPI", helloHandler)
	go func() {
		http.ListenAndServe(":"+config.SecondPort, nil)
	}()
}

// Resp is json model
type Resp struct {
	Message string
	Error   string
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	go createRmqConnect()

	fmt.Println("1 zapros")
	w.Header().Set("Content-Type", "application/json")
	r.ParseForm()
	//value := "usering:1234@" + config.Bridge + ":3306" + "/app_stats"
	value := config.SQLConnectLink
	//value := "usering:1234" + "@" + config.Bridge + ":3306" + "/app_stats"
	fmt.Println(value)

	db, err := sql.Open("mysql", value)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(db.Ping())
	}

	defer db.Close()

	resp := Resp{
		Message: "hello",
	}
	respJ, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respJ)

}
