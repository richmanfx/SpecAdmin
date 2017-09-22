package modules

import (
	"github.com/gin-gonic/gin"
	"./handlers"
	"./auth"
)

func InitRoutes(router *gin.Engine) {

		// Роутинг страницы: метод, путь -> обработчик
		router.Handle("GET", "/spec-admin", auth.AuthRequired(), handlers.ShowIndexPage)
		router.Handle("POST", "/spec-admin/show-suites", auth.AuthRequired(), handlers.ShowSuitesIndex)

		router.Handle("POST", "/spec-admin/add-group", auth.AuthRequired(), handlers.AddGroup)
		router.Handle("POST", "/spec-admin/del-group", auth.AuthRequired(), handlers.DelGroup)
		router.Handle("POST", "/spec-admin/edit-group", auth.AuthRequired(), handlers.EditGroup)

		router.Handle("POST", "/spec-admin/add-suite", auth.AuthRequired(), handlers.AddSuite)
		router.Handle("POST", "/spec-admin/del-suite", auth.AuthRequired(), handlers.DelSuite)
		router.Handle("POST", "/spec-admin/edit-suite", auth.AuthRequired(), handlers.EditSuite)
		router.Handle("POST", "/spec-admin/rename-suite", auth.AuthRequired(), handlers.RenameSuite)
		router.Handle("POST", "/spec-admin/update-after-edit-suite", auth.AuthRequired(), handlers.UpdateAfterEditSuite)
		router.Handle("POST", "/spec-admin/print-suites-scripts", auth.AuthRequired(), handlers.CreateScriptsPdf) // for AJAX


		router.Handle("POST", "/spec-admin/add-script", auth.AuthRequired(), handlers.AddScript)
		router.Handle("POST", "/spec-admin/del-script", auth.AuthRequired(), handlers.DelScript)
		router.Handle("POST", "/spec-admin/edit-script", auth.AuthRequired(), handlers.EditScript)
		router.Handle("POST", "/spec-admin/update-after-edit-script", auth.AuthRequired(), handlers.UpdateAfterEditScript)
		router.Handle("POST", "/spec-admin/print-scripts-steps", auth.AuthRequired(), handlers.CreateStepsPdf) // for AJAX

		router.Handle("POST", "/spec-admin/add-step", auth.AuthRequired(), handlers.AddStep)
		router.Handle("POST", "/spec-admin/del-step", auth.AuthRequired(), handlers.DelStep)
		router.Handle("POST", "/spec-admin/edit-step", auth.AuthRequired(), handlers.EditStep)
		router.Handle("POST", "/spec-admin/update-after-edit-step", auth.AuthRequired(), handlers.UpdateAfterEditStep)
		router.Handle("POST", "/spec-admin/get-steps-options", auth.AuthRequired(), handlers.GetStepsOptions)     // for AJAX
		router.Handle("POST", "/spec-admin/del-screen-shot", auth.AuthRequired(), handlers.DelScreenShotFromStep) // for AJAX
		router.Handle("POST", "/spec-admin/copy-step-in-clipboard", auth.AuthRequired(), handlers.CopyStepInClipboard) // for AJAX

		router.Handle("GET", "/spec-admin/edit-config", auth.AuthRequired(), handlers.EditConfig)
		router.Handle("POST", "/spec-admin/save-config", auth.AuthRequired(), handlers.SaveConfig)
		router.Handle("GET", "/spec-admin/del-old-screenshots-file", auth.AuthRequired(), handlers.DelUnusedScreenShotsFile)

		router.Handle("GET", "/spec-admin/login", auth.Login)
		router.Handle("POST", "/spec-admin/login-processing", auth.Authorization)
		router.Handle("GET", "/spec-admin/users-config", auth.AuthRequired(), handlers.UsersConfig)
		router.Handle("POST", "/spec-admin/create-user", auth.AuthRequired(), handlers.CreateUser)
		router.Handle("POST", "/spec-admin/delete-user", auth.AuthRequired(), handlers.DeleteUser)
		router.Handle("POST", "/spec-admin/save-user", auth.AuthRequired(), handlers.SaveUser)
		router.Handle("POST", "/spec-admin/change-password", auth.AuthRequired(), auth.ChangePassword)
		router.Handle("POST", "/spec-admin/logout", auth.Logout)
	}
