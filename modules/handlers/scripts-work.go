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
	scriptsSuite := context.PostForm("script_suite")

	// Группа Сюиты
	suite, err := helpers.GetSuite(scriptsSuite)
	suitesGroup := suite.Group

	err = helpers.AddTestScript(newScript, scriptSerialNumber, scriptsSuite)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	fmt.Sprintf("Ошибка при добавлении сценария '%s' в сюиту '%s'.", newScript, scriptsSuite),
				"message3": 	fmt.Sprintf("%s: ", err),
				"SuitesGroup":	suitesGroup,
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Информация",
				"message1": 	fmt.Sprintf("Сценарий '%s' успешно добавлен в сюиту '%s'.", newScript, scriptsSuite),
				"message2": 	"",
				"message3": 	"",
				"SuitesGroup":	suitesGroup,
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}

// Удалить сценарий из БД
func DelScript(context *gin.Context)  {
	helpers.SetLogFormat()

	// Данные из формы
	deletedScript := context.PostForm("script")
	scriptsSuite := context.PostForm("scripts_suite")

	err := helpers.DelTestScript(deletedScript, scriptsSuite)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	fmt.Sprintf("Ошибка при удалении скрипта '%s'.", deletedScript),
				"message3": 	fmt.Sprintf("%s: ", err),
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Информация",
				"message1": 	fmt.Sprintf("Скрипт '%s' успешно удалён.", deletedScript),
				"message2": 	"",	"message3": "",
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}


// Редактировать Сценарий, данные получить из БД
func EditScript(context *gin.Context)  {
	helpers.SetLogFormat()

	// Данные из формы
	editedScript := context.PostForm("script")
	scriptsSuite := context.PostForm("scripts_suite")

	log.Debugf("Редактируется Сценарий '%v' в Сюите '%v'", editedScript, scriptsSuite)

	// Получить данные о сценарии из БД
	var script models.Script
	var scriptId int
	var err error
	script, scriptId, err = helpers.GetScript(editedScript, scriptsSuite)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	fmt.Sprintf("Ошибка получения данных о сценарии '%s'.", editedScript),
				"message3": 	fmt.Sprintf("%s: ", err),
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		// Вывести данные для редактирования
		context.HTML(http.StatusOK, "edit-script.html",
			gin.H{
				"title": 		"Редактирование сценария",
				"script": 		script,
				"scriptId":		scriptId,
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}

// Обновить в БД Сценарий после редактирования
func UpdateAfterEditScript(context *gin.Context) {
	helpers.SetLogFormat()
	var err error
	var suite models.Suite
	var suitesGroup string
	var scriptName string
	var scriptSerialNumber string
	var scriptsSuite string

	//Данные из формы
	scriptId, err := strconv.Atoi(context.PostForm("hidden_id"))
	if err == nil {

		scriptName = context.PostForm("script")
		scriptSerialNumber = context.PostForm("scripts_serial_number")
		scriptsSuite = context.PostForm("scripts_suite")

		// Группа Сюиты
		suite, err = helpers.GetSuite(scriptsSuite)
		suitesGroup = suite.Group

		// Обновить в БД
		err = helpers.UpdateTestScript(scriptId, scriptName, scriptSerialNumber, scriptsSuite)
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	fmt.Sprintf("Ошибка при обновлении сценария '%s' в сюите '%s'.", scriptName, scriptsSuite),
				"message3": 	fmt.Sprintf("%s: ", err),
				"SuitesGroup":	suitesGroup,
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Информация",
				"message1": 	fmt.Sprintf("Сценарий '%s' успешно обновлён.", scriptName),
				"message2": 	"",
				"message3": 	"",
				"SuitesGroup":	suitesGroup,
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}


// Вывод для печати Шагов Сценария по имени Сценария и имени его Сюиты
func CreateStepsPdf(context *gin.Context)  {

	helpers.SetLogFormat()

	// Данные из AJAX запроса
	scriptName := context.PostForm("scriptName")
	scriptsSuite := context.PostForm("suiteName")
	log.Debugf("Данные из POST запроса AJAX: '%v' и '%v'", scriptName, scriptsSuite)

	// Получить из базы id сценария по имени Сценария и имени его Сюиты
	_, scriptId, err := helpers.GetScript(scriptName, scriptsSuite)

	// Получить Шаги из БД только для заданных по ID Сценариев
	scriptsIdList := append(make([]int, 0, 1), scriptId)		// Слайс только из одного Id
	stepsList, err := helpers.GetStepsListForSpecifiedScripts(scriptsIdList)
	log.Debugf("%v - %v", stepsList, err)

	// Сгенерировать PDF
	err = helpers.GetScripsStepsPdf(scriptsSuite, scriptName, stepsList)
	if err != nil {
		context.AbortWithStatus(404)
	}

	if err == nil {
		context.Abort()
		context.Redirect(http.StatusSeeOther, "ok")
	} else {
		log.Errorf("Ошибка при генерации PDF-файла с шагами сценария: %v", err)
		result := gin.H{"scriptId": err}
		context.JSON(http.StatusOK, result)
	}
}
