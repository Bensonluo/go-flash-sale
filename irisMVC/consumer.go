package main

import (
	"fmt"
	"go-flash-sale/irisMVC/common"
	"go-flash-sale/irisMVC/rabbitmq"
	"go-flash-sale/irisMVC/repositories"
	"go-flash-sale/irisMVC/service"
)

func main() {
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}

	//create product
	product := repositories.NewProductManager("product", db)
	//create product service
	productService := service.NewProductService(product)

	order := repositories.NewOrderManagerRepository("order", db)
	orderService := service.NewOrderService(order)

	rabbitmqConsumeSimple := rabbitmq.NewRabbitMQSimple("goFlashSale")
	rabbitmqConsumeSimple.ConsumeSimple(orderService, productService)
}
