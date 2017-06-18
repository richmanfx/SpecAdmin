package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"fmt"
	"net/http"
)

// Добавить в БД новый сценарий
func AddScript(context *gin.Context)  {

	helpers.SetLogFormat()

	// Данные из формы
	newScript := context.PostForm("script")
	scriptSerialNumber := context.PostForm("scripts_serial_number")
	scriptSuite := context.PostForm("script_suite")

	err := helpers.AddTestScript(newScript, scriptSerialNumber, scriptSuite)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при добавлении сценария '%s' в сюиту '%s'.", newScript, scriptSuite),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Сценарий '%s' успешно добавлен в сюиту '%s'.", newScript, scriptSuite),
				"message2": "",
				"message3": "",
			},
		)
	}
}
