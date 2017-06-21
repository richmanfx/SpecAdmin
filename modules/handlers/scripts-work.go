package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"fmt"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"../../models"
	"strconv"
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

// Удалить сценарий из БД
func DelScript(context *gin.Context)  {
	helpers.SetLogFormat()

	deletedScript := context.PostForm("script")			// Скрипт из формы
	err := helpers.DelTestScript(deletedScript)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при удалении скрипта '%s'.", deletedScript),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Скрипт '%s' успешно удалён.", deletedScript),
				"message2": "",	"message3": "",
			},
		)
	}
}


// Редактировать Сценарий, данные получить из БД
func EditScript(context *gin.Context)  {
	helpers.SetLogFormat()

	editedScript := context.PostForm("script")					// Скрипт из формы
	log.Infof("Редактируется Сценарий: '%v'.", editedScript)

	// Получить данные о сценарии из БД
	var script models.Script
	var scriptId int
	var err error
	script, scriptId, err = helpers.GetScript(editedScript)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка получения данных о сценарии '%s'.", editedScript),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		// Вывести данные для редактирования
		context.HTML(http.StatusOK, "edit-script.html",
			gin.H{
				"title": 	"Редактирование сценария",
				"Version":	Version,
				"script": 	script,
				"scriptId":	scriptId,
			},
		)
	}
}

// Обновить в БД скрипт после редактирования
func UpdateAfterEditScript(context *gin.Context) {
	helpers.SetLogFormat()

	//Данные из формы
	scriptId, err := strconv.Atoi(context.PostForm("hidden_id"))
	if err != nil {
		panic(err)
	}

	scriptName := context.PostForm("script")
	scriptSerialNumber := context.PostForm("scripts_serial_number")
	scriptsSuite := context.PostForm("scripts_suite")

	// Обновить в БД
	err = helpers.UpdateTestScript(scriptId, scriptName, scriptSerialNumber, scriptsSuite)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при обновлении сценария '%s' в сюите '%s'.", scriptName, scriptsSuite),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Сценарий '%s' успешно обновлён.", scriptName),
				"message2": "",
				"message3": "",
			},
		)
	}
}
















