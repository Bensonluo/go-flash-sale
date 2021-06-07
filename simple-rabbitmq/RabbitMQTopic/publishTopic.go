package main

import (
	"fmt"
	"go-flash-sale/simple-rabbitmq/RabbitMQ"
	"strconv"
	"time"
)

func main()  {
	bensonOne:=RabbitMQ.NewRabbitMQTopic("exChangeTopic","benson.topic.one")
	bensonTwo:=RabbitMQ.NewRabbitMQTopic("exChangeTopic","benson.topic.two")
	for i := 0; i <= 10; i++ {
		bensonOne.PublishTopic("Hello benson topic one!" + strconv.Itoa(i))
		bensonTwo.PublishTopic("Hello benson topic Two!" + strconv.Itoa(i))
		time.Sleep(1 * time.Second)
		fmt.Println(i)
	}
	
}
