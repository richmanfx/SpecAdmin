package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"fmt"
	"net/http"
)

// Добавить группу в базу
func AddGroup(context *gin.Context)  {

	newGroup := context.PostForm("group")		// Группа из формы
	err := helpers.AddTestGroup(newGroup)

	if err != nil {
		// Ошибка при добавлении группы
		context.HTML(http.StatusOK, "message-modal.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при добавлении группы тестов '%s'.", newGroup),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message-modal.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Группа тестов '%s' добавлена успешно.", newGroup),
				"message2": "",
				"message3": "",
			},
		)
	}
}

func DelGroup(context *gin.Context)  {

	deletedGroup := context.PostForm("group")		// Группа из формы
	err := helpers.DelTestGroup(deletedGroup)

	if err != nil {
		// Ошибка при удалении группы
		context.HTML(http.StatusOK, "message-modal.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при удалении группы тестов '%s'.", deletedGroup),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message-modal.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Группа тестов '%s' удалена успешно.", deletedGroup),
				"message2": "",
				"message3": "",
			},
		)
	}
}

func EditGroup(context *gin.Context)  {

	// Прежнее и новое имя для Группы из формы
	oldGroup := context.PostForm("old_group")
	newGroup := context.PostForm("new_group")
	err := helpers.EditTestGroup(oldGroup, newGroup)

	if err != nil {
		// Ошибка при редактировании группы
		context.HTML(http.StatusOK, "message-modal.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при редактировании группы тестов '%s'.", oldGroup),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message-modal.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf(
					"Группа тестов '%s' успешно изменена на '%s'.",
					oldGroup,
					newGroup),
				"message2": "",
				"message3": "",
			},
		)
	}
}