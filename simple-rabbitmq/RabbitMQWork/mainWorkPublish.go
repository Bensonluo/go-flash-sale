package main

import (
	"fmt"
	"go-flash-sale/simple-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQSimple("" +
		"SimpleQueue")

	for i := 0; i <= 100; i++ {
		rabbitmq.PublishSimple("Hello world!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
}
