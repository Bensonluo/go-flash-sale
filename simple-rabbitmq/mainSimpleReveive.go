package main

import (
	"go-flash-sale/simple-rabbitmq/RabbitMQ"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("simpleQueue")
	rabbitmq.ConsumeSimple()
}
