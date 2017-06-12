package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"fmt"
	"net/http"
)

func AddGroup(context *gin.Context)  {

	// Добавить группу в базу
	newGroup := context.PostForm("group")		// Группа из формы

	err := helpers.AddTestGroup(newGroup)

	if err != nil {
		// Ошибка при добавлении группы
		context.HTML(
			http.StatusOK,
			"message-modal.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при добавлении группы тестов '%s'.", newGroup),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		// Вывести сообщение - группа добавлена, ссылка возврата на главную страницу
		context.HTML(
			http.StatusOK,
			"message-modal.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Группа тестов '%s' добавлена успешно.", newGroup),
				"message2": "",
				"message3": "",
			},
		)
	}
}
