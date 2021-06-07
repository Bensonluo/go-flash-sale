package main

import (
	"fmt"
	"go-flash-sale/simple-rabbitmq/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("simpleQueue")
	rabbitmq.PublishSimple("Hello world!")
	fmt.Println("sent")
}
