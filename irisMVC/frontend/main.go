package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"go-flash-sale/irisMVC/common"
	"go-flash-sale/irisMVC/frontend/web/controllers"
	"go-flash-sale/irisMVC/repositories"
	"go-flash-sale/irisMVC/service"

	"github.com/kataras/iris/v12/sessions"
	"time"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")

	template := iris.HTML("./irisMVC/frontend/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	app.HandleDir("/public", iris.Dir("./irisMVC/frontend/web/public"))
	app.HandleDir("/html", iris.Dir("./irisMVC/frontend/web/htmlProductShow"))

	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})

	db, err := common.NewMysqlConn()
	if err != nil {

	}
	sess := sessions.New(sessions.Config{
		Cookie:"AdminCookie",
		Expires:600*time.Minute,
	})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserRepository("user", db)
	userService := service.NewService(user)
	userPro := mvc.New(app.Party("/user"))
	userPro.Register(userService, ctx,sess.Start)
	userPro.Handle(new(controllers.UserController))

	app.Run(
		iris.Addr("0.0.0.0:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}