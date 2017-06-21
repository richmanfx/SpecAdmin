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
	router.Handle("POST", "/spec-admin/edit-group", handlers.EditGroup)

	router.Handle("POST", "/spec-admin/add-suite", handlers.AddSuite)
	router.Handle("POST", "/spec-admin/del-suite", handlers.DelSuite)
	router.Handle("POST", "/spec-admin/edit-suite", handlers.EditSuite)
	router.Handle("POST", "/spec-admin/update-after-edit-suite", handlers.UpdateAfterEditSuite)

	router.Handle("POST", "/spec-admin/add-script", handlers.AddScript)
	router.Handle("POST", "/spec-admin/del-script", handlers.DelScript)
	router.Handle("POST", "/spec-admin/edit-script", handlers.EditScript)
	router.Handle("POST", "/spec-admin/update-after-edit-script", handlers.UpdateAfterEditScript)

}

