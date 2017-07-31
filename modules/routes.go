package modules

import (
	"github.com/gin-gonic/gin"
	"./handlers"
	"./auth"
)

func InitRoutes(router *gin.Engine) {

	// Роутинг страницы: метод, путь -> обработчик
	router.Handle("GET", "/spec-admin", auth.AuthRequired(), handlers.ShowIndexPage)	// Навесил мидлеварю
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
	router.Handle("POST", "/spec-admin/del-screen-shot", handlers.DelScreenShotFromStep)		// for AJAX

	router.Handle("GET", "/spec-admin/edit-config", handlers.EditConfig)
	router.Handle("POST", "/spec-admin/save-config", handlers.SaveConfig)
	router.Handle("GET", "/spec-admin/del-old-screenshots-file", handlers.DelUnusedScreenShotsFile)

	router.Handle("GET", "/spec-admin/login", auth.Login)
	router.Handle("POST", "/spec-admin/login-processing", auth.Authorization)
	router.Handle("GET", "/spec-admin/users-config", handlers.UsersConfig)
	router.Handle("POST", "/spec-admin/create-user", handlers.CreateUser)
	router.Handle("POST", "/spec-admin/delete-user", handlers.DeleteUser)
	router.Handle("POST", "/spec-admin/change-password", auth.ChangePassword)
	router.Handle("POST", "/spec-admin/logout", auth.Logout)


	//authorized := router.Group("/spec-admin/")
	//authorized.Use(gin.BasicAuth(gin.Accounts{
	//	"user1": "user1", // user:user1 password:user1
	//	"user2": "user2", // user:user2 password:user2
	//}))
}
