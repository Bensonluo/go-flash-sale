package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/opentracing/opentracing-go/log"
	"go-flash-sale/irisMVC/backend/web/controllers"
	"go-flash-sale/irisMVC/common"
	"go-flash-sale/irisMVC/repositories"
	"go-flash-sale/irisMVC/service"
)


func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	template := iris.HTML("./backend/web/views", ".html").Layout(
		"shared/layout.html").Reload(
			true)
	app.RegisterView(template)

	app.HandleDir("/assets", iris.Dir("./backend/web/assets"))

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "Unknown Error has occurred!"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	//connect db
	db, err := common.NewMysqlConn()
	if err != nil {
		log.Error(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//controller
	productRepository := repositories.NewProductManager("product", db)
	productService := service.NewProductService(productRepository)
	productParty := app.Party("/product")
	product := mvc.New(productParty)
	product.Register(ctx, productService)
	product.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations)
}


