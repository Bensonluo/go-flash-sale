package main

import "go-flash-sale/simple-rabbitmq/RabbitMQ"

func main()  {
	bensonOne:=RabbitMQ.NewRabbitMQTopic("exChangeTopic","#")
	bensonOne.RecieveTopic()
}
