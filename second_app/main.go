package main

import (
	"MyProjects/go_app/models/appConfig"
	"MyProjects/go_app/models/dbModels"
	"MyProjects/go_app/servises/rmqServise"
	"MyProjects/go_app/servises/sqlServise"
	"encoding/json"
	"fmt"
	"net/http"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var config = appConfig.Get()
var rmq = rmqServise.Connect(config.RMQConnectLink)
var sqlC = sqlServise.Connect(config.SQLConnectLink)

func main() {
	fmt.Println("Start")
	createHelloHandler()
	rmq.Handle(handler)

	c1 := make(chan string)
	<-c1
}

func handler(body []byte) bool {
	addTask := dbModels.User{}
	err := json.Unmarshal(body, &addTask)
	if err != nil {
		fmt.Println(err)
	}
	sqlC.Exec(
		"insert into productdb.Products (model, company, price) values (?, ?, ?)",
		"iPhone X", "Apple", 72000,
	)
	return true
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
	w.Header().Set("Content-Type", "application/json")

	resp := Resp{
		Message: "hello",
	}
	respJ, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusOK)
	w.Write(respJ)
}

func sqlConnect() {
	db, err := sql.Open("mysql", config.SQLConnectLink)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(db.Ping())
	}
	defer db.Close()
}
