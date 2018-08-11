package main

import (
	"./modules"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func main() {

	// Режим работы gin - на продакшене делать "ReleaseMode"
	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	// Роутер по-умолчанию в Gin
	//Router = gin.Default()
	Router = gin.New()

	store := cookie.NewStore([]byte("secret"))
	Router.Use(sessions.Sessions("mysession", store))

	// Загрузить шаблоны
	Router.LoadHTMLGlob("templates/*")

	// Загрузить статику
	Router.Static("/css", "css")
	Router.Static("/js", "js")
	Router.Static("/images", "images")
	Router.Static("/pdf", "pdf")

	// Проинитить роуты
	modules.InitRoutes(Router)

	// Запустить приложение
	Router.Run(":9094")
}
