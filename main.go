package main

import (
	"SpecAdmin/modules"
	"fmt"
	log "github.com/Sirupsen/logrus"
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
	tcpPort := "3010"
	err := Router.Run(fmt.Sprintf(":%s", tcpPort))
	if err != nil {
		log.Errorf("Ошибка запуска приложение на порту '%s'", tcpPort)
	}
}
