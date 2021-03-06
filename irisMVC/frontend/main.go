package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go-flash-sale/irisMVC/common"
	"go-flash-sale/irisMVC/frontend/middleware"
	"go-flash-sale/irisMVC/frontend/web/controllers"
	"go-flash-sale/irisMVC/rabbitmq"
	"go-flash-sale/irisMVC/repositories"
	"go-flash-sale/irisMVC/service"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	template := iris.HTML("./irisMVC/frontend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	//app.StaticWeb
	app.HandleDir("/public", iris.Dir("./irisMVC/frontend/web/public"))
	app.HandleDir("/html", iris.Dir("./irisMVC/frontend/web/htmlProductShow"))

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "errors happen！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {

	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserRepository("user", db)
	userService := service.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService)
	userPro.Handle(new(controllers.UserController))

	rabbitMQ := rabbitmq.NewRabbitMQSimple("goFlashSale")

	order := repositories.NewOrderManagerRepository("order", db)
	orderService := service.NewOrderService(order)

	product := repositories.NewProductManager("product", db)
	productService := service.NewProductService(product)
	proProduct := app.Party("/product")
	pro := mvc.New(proProduct)
	proProduct.Use(middleware.AuthConProduct)

	pro.Register(productService, orderService, ctx, rabbitMQ)
	pro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
