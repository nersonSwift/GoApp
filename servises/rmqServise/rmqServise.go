package rmqServise

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type RMQConnection struct {
	Link    string
	Handler func([]byte) bool

	Connection *amqp.Connection
	Channel    *amqp.Channel
	Queue      amqp.Queue

	Cash []interface{}
}

func Connect(link string) *RMQConnection {
	rmqConnection := RMQConnection{
		Link: link,
		Cash: make([]interface{}, 0),
	}
	rmqConnection.Connect()
	return &rmqConnection
}

func (rmq *RMQConnection) Connect() {
	rmq.SetConnect()
	rmq.SetChannel()
	rmq.SetQueue()
	rmq.SetHandler()
	rmq.SetErrConnectHandler()

	go func() {
		if rmq.Connection == nil {
			timer := time.NewTimer(1 * time.Second)
			<-timer.C
			if rmq.Connection == nil {
				rmq.Reconnect()
			}
		} else {
			cash := rmq.Cash
			rmq.Cash = make([]interface{}, 0)
			for _, body := range cash {
				fmt.Println("body:", body)
				rmq.SendMassage(body)
			}
		}
	}()
}

func (rmq *RMQConnection) Reconnect() {
	fmt.Println("Reconnect")
	rmq.Connection = nil
	rmq.Channel = nil
	rmq.Connect()
}

func (rmq *RMQConnection) SetConnect() {
	if rmq.Connection != nil {
		return
	}
	connect, err := amqp.Dial(rmq.Link)
	if err != nil {
		fmt.Println(err)
		return
	}
	rmq.Connection = connect
}

func (rmq *RMQConnection) SetChannel() {
	if rmq.Channel != nil || rmq.Connection == nil {
		return
	}

	channel, err := rmq.Connection.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	rmq.Channel = channel
}

func (rmq *RMQConnection) SetErrConnectHandler() {
	if rmq.Connection == nil {
		return
	}
	go func() {
		<-rmq.Connection.NotifyClose(make(chan *amqp.Error)) //Listen to NotifyClose
		fmt.Println("ErrConnectHandler")
		rmq.Reconnect()
	}()
}

func (rmq *RMQConnection) SetQueue() {
	if rmq.Channel == nil {
		return
	}

	queue, err := rmq.Channel.QueueDeclare("add", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	rmq.Queue = queue
}

func (rmq *RMQConnection) SendMassage(body interface{}) {
	if rmq.Connection == nil || rmq.Channel == nil {
		rmq.Cash = append(rmq.Cash, body)
		return
	}

	b, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = rmq.Channel.Publish("", rmq.Queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         b,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (rmq *RMQConnection) Handle(handler func([]byte) bool) {
	rmq.Handler = handler
	rmq.SetHandler()
}

func (rmq *RMQConnection) SetHandler() {
	if rmq.Channel == nil || rmq.Handler == nil {
		return
	}
	err := rmq.Channel.Qos(1, 0, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	messageChannel, err := rmq.Channel.Consume(
		rmq.Queue.Name,
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

	go func() {
		for d := range messageChannel {
			if err := d.Ack(!rmq.Handler(d.Body)); err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("End")
	}()
}

// func createRmqConnect1() {
// 	err = channel.Qos(1, 0, false)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	messageChannel, err := channel.Consume(
// 		queue.Name,
// 		"",
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
//
// 	stopChan := make(chan bool)
// 	go func() {
// 		for d := range messageChannel {
// 			addTask := Mess{}
// 			err := json.Unmarshal(d.Body, &addTask)
// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			fmt.Println("Mess: " + addTask.Message)
//
// 			if err := d.Ack(false); err != nil {
// 				log.Printf("Error acknowledging message : %s", err)
// 			} else {
// 				log.Printf("Acknowledged message")
// 			}
// 		}
// 		fmt.Println("End")
// 	}()
// 	<-stopChan
// }
//
// func (c *Config) Set() {
// 	var config Config
// 	err := envconfig.Process("", &config)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	*c = config
// }
//
// // Get...
// func Get() Config {
// 	var config Config
// 	config.Set()
// 	return config
// }
