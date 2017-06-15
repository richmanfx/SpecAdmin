package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"fmt"
	"net/http"
)

// Добавить сюиту в базу
func AddSuite(context *gin.Context)  {

	newSuite := context.PostForm("suite")							// Сюита из формы
	suitesDescription := context.PostForm("suites_description")		// Описание Сюиты
	suitesGroup := context.PostForm("suites_group")					// Группа Сюиты
	err := helpers.AddTestSuite(newSuite, suitesDescription, suitesGroup)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при добавлении сюиты '%s' в группу '%s'.", newSuite, suitesGroup),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Сюита '%s' успешно добавлена в группу '%s'.", newSuite, suitesGroup),
				"message2": "",
				"message3": "",
			},
		)
	}
}

// Удалить сюиту из базы
func DelSuite(context *gin.Context)  {

	deletedSuite := context.PostForm("suite")							// Сюита из формы
	err := helpers.DelTestSuite(deletedSuite)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при удалении сюиты '%s'.", deletedSuite),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Сюита '%s' удалена успешно.", deletedSuite),
				"message2": "",	"message3": "",
			},
		)
	}
}