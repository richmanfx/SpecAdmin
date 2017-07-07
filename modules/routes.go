package modules

import (
	"github.com/gin-gonic/gin"
	"./handlers"
)

func InitRoutes(router *gin.Engine) {

	// Роутинг страницы: метод, путь -> обработчик
	router.Handle("GET", "/spec-admin", handlers.ShowIndexPage)
	router.Handle("POST", "/spec-admin/show-suites", handlers.ShowSuitesIndex)

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

	router.Handle("POST", "/spec-admin/add-step", handlers.AddStep)
	router.Handle("POST", "/spec-admin/del-step", handlers.DelStep)
	router.Handle("POST", "/spec-admin/edit-step", handlers.EditStep)
	router.Handle("POST", "/spec-admin/update-after-edit-step", handlers.UpdateAfterEditStep)
	router.Handle("POST", "/spec-admin/get-steps-options", handlers.GetStepsOptions)			// for AJAX


	router.Handle("GET", "/spec-admin/edit-config", handlers.EditConfig)
	router.Handle("POST", "/spec-admin/save-config", handlers.SaveConfig)
}
