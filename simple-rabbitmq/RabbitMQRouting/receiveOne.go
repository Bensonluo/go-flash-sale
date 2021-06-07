package main

import "go-flash-sale/simple-rabbitmq/RabbitMQ"

func main()  {
	imoocOne:=RabbitMQ.NewRabbitMQRouting("exMessage","key_one")
	imoocOne.ReceiveRouting()
}
