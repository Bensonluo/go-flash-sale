package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go-flash-sale/irisMVC/service"
)

type OrderController struct {
	Ctx iris.Context
	OrderService service.IOrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray,err:=o.OrderService.GetAllOrderInfo()
	if err !=nil {
		o.Ctx.Application().Logger().Debug("Find Orders Fails")
	}

	return mvc.View{
		Name:"order/view.html",
		Data:iris.Map{
			"order":orderArray,
		},
	}
}