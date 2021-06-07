package main

import (
	"fmt"
	"go-flash-sale/simple-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	rabbitmq := RabbitMQ.NewRabbitMQPubSub("" +
		"newProduct")
	for i := 0; i < 100; i++ {
		rabbitmq.PublishPub("Produce the" +
			strconv.Itoa(i) + "th" + "data")
		fmt.Println("Produce the" +
			strconv.Itoa(i) + "th" + "data")
		time.Sleep(1 * time.Second)
	}

}
