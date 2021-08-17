package main

import (
	"fmt"
	"go-flash-sale/simple-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main() {
	imoocOne := RabbitMQ.NewRabbitMQRouting("exMessage", "key_one")
	imoocTwo := RabbitMQ.NewRabbitMQRouting("exMessage", "key_two")
	for i := 0; i <= 10; i++ {
		imoocOne.PublishRouting("Hello World one!" + strconv.Itoa(i))
		imoocTwo.PublishRouting("Hello World Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}

}
