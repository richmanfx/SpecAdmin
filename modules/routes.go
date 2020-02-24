package modules

import (
	"SpecAdmin/modules/auth"
	"SpecAdmin/modules/handlers"
	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	// Роутинг страницы: метод, путь -> обработчик
	router.Handle("GET", "/spec-admin", auth.AuthorizationRequired(), handlers.ShowIndexPage)
	router.Handle("POST", "/spec-admin/show-suites", auth.AuthorizationRequired(), handlers.ShowSuitesIndex)

	router.Handle("POST", "/spec-admin/add-group", auth.AuthorizationRequired(), handlers.AddGroup)
	router.Handle("POST", "/spec-admin/del-group", auth.AuthorizationRequired(), handlers.DelGroup)
	router.Handle("POST", "/spec-admin/edit-group", auth.AuthorizationRequired(), handlers.EditGroup)

	router.Handle("POST", "/spec-admin/add-suite", auth.AuthorizationRequired(), handlers.AddSuite)
	router.Handle("POST", "/spec-admin/del-suite", auth.AuthorizationRequired(), handlers.DelSuite)
	router.Handle("POST", "/spec-admin/edit-suite", auth.AuthorizationRequired(), handlers.EditSuite)
	router.Handle("POST", "/spec-admin/rename-suite", auth.AuthorizationRequired(), handlers.RenameSuite)
	router.Handle("POST", "/spec-admin/update-after-edit-suite", auth.AuthorizationRequired(), handlers.UpdateAfterEditSuite)
	router.Handle("POST", "/spec-admin/print-suites-scripts", auth.AuthorizationRequired(), handlers.CreateScriptsPdf) // for AJAX

	router.Handle("POST", "/spec-admin/add-script", auth.AuthorizationRequired(), handlers.AddScript)
	router.Handle("POST", "/spec-admin/del-script", auth.AuthorizationRequired(), handlers.DelScript)
	router.Handle("POST", "/spec-admin/edit-script", auth.AuthorizationRequired(), handlers.EditScript)
	router.Handle("POST", "/spec-admin/update-after-edit-script", auth.AuthorizationRequired(), handlers.UpdateAfterEditScript)
	router.Handle("POST", "/spec-admin/print-scripts-steps", auth.AuthorizationRequired(), handlers.CreateStepsPdf) // for AJAX

	router.Handle("POST", "/spec-admin/add-step", auth.AuthorizationRequired(), handlers.AddStep)
	router.Handle("POST", "/spec-admin/del-step", auth.AuthorizationRequired(), handlers.DelStep)
	router.Handle("POST", "/spec-admin/edit-step", auth.AuthorizationRequired(), handlers.EditStep)
	router.Handle("POST", "/spec-admin/update-after-edit-step", auth.AuthorizationRequired(), handlers.UpdateAfterEditStep)
	router.Handle("POST", "/spec-admin/get-steps-options", auth.AuthorizationRequired(), handlers.GetStepsOptions)          // for AJAX
	router.Handle("POST", "/spec-admin/del-screen-shot", auth.AuthorizationRequired(), handlers.DelScreenShotFromStep)      // for AJAX
	router.Handle("POST", "/spec-admin/copy-step-in-clipboard", auth.AuthorizationRequired(), handlers.CopyStepInClipboard) // for AJAX
	router.Handle("POST", "/spec-admin/get-step-from-buffer", auth.AuthorizationRequired(), handlers.GetStepFromBuffer)     // for AJAX

	router.Handle("GET", "/spec-admin/edit-config", auth.AuthorizationRequired(), handlers.EditConfig)
	router.Handle("POST", "/spec-admin/save-config", auth.AuthorizationRequired(), handlers.SaveConfig)
	router.Handle("GET", "/spec-admin/del-old-screenshots-file", auth.AuthorizationRequired(), handlers.DelUnusedScreenShotsFile)

	router.Handle("GET", "/spec-admin/login", auth.Login)
	router.Handle("POST", "/spec-admin/login-processing", auth.Authorization)
	router.Handle("GET", "/spec-admin/users-config", auth.AuthorizationRequired(), handlers.UsersConfig)
	router.Handle("POST", "/spec-admin/create-user", auth.AuthorizationRequired(), handlers.CreateUser)
	router.Handle("POST", "/spec-admin/delete-user", auth.AuthorizationRequired(), handlers.DeleteUser)
	router.Handle("POST", "/spec-admin/save-user", auth.AuthorizationRequired(), handlers.SaveUser)
	router.Handle("POST", "/spec-admin/change-password", auth.AuthorizationRequired(), auth.ChangePassword)
	router.Handle("POST", "/spec-admin/logout", auth.Logout)
}
