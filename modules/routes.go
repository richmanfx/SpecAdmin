package modules

import (
	"github.com/gin-gonic/gin"
	"./handlers"
)

func InitRoutes(router *gin.Engine) {

	// Роутинг страницы: метод, путь -> обработчик
	router.Handle("GET", "/spec-admin", handlers.ShowIndexPage)
	router.Handle("POST", "/spec-admin/add-group", handlers.AddGroup)
	router.Handle("POST", "/spec-admin/del-group", handlers.DelGroup)
}

