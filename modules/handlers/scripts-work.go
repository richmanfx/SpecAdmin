package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"fmt"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"../../models"
	"strconv"
	"path/filepath"
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

	helpers.CloseConnectToDB()

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при добавлении сценария '%s' в сюиту '%s'.", newScript, scriptsSuite),
				"message3": fmt.Sprintf("%s: ", err),
				"SuitesGroup":	suitesGroup,
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Сценарий '%s' успешно добавлен в сюиту '%s'.", newScript, scriptsSuite),
				"message2": "",
				"message3": "",
				"SuitesGroup":	suitesGroup,
				"Version":	Version,
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

	helpers.CloseConnectToDB()

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при удалении скрипта '%s'.", deletedScript),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Скрипт '%s' успешно удалён.", deletedScript),
				"message2": "",	"message3": "",
				"Version":	Version,
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

	helpers.CloseConnectToDB()

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка получения данных о сценарии '%s'.", editedScript),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		// Вывести данные для редактирования
		context.HTML(http.StatusOK, "edit-script.html",
			gin.H{
				"title": 	"Редактирование сценария",
				"script": 	script,
				"scriptId":	scriptId,
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}

// Обновить в БД скрипт после редактирования
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

	helpers.CloseConnectToDB()

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
func GetScriptsId(context *gin.Context)  {

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

	// Закрыть соединение с БД
	helpers.CloseConnectToDB()

	// Сгенерировать PDF
	pdfFileName, err := helpers.GetScripsStepsPdf(scriptsSuite, scriptName, stepsList)
	rootDir := "C:\\Users\\Admin\\GoglandProjects\\SpecAdmin\\"		// TODO: В директорию к скриншотам оформить файл
	//name := context.Param(pdfFileName)
	filePath, err :=  filepath.Abs(rootDir + pdfFileName)
	if err != nil {
		context.AbortWithStatus(404)
	}

	if err == nil {
		//result := gin.H{"scriptId": scriptId}
		//context.JSON(http.StatusOK, result)
		context.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
		context.Header("Content-Type", "application/pdf")
		context.Header("Content-Transfer-Encoding", "binary")
		context.Header("Content-Length", strconv.Itoa(len(filePath)))
		context.Header("Content-Disposition", fmt.Sprintf("inline; filename='%s'", pdfFileName))

		context.File(filePath)		// Отправить PDF-файл
	} else {
		log.Errorf("Ошибка при генерации PDF: %v", err)
		result := gin.H{"scriptId": err}
		context.JSON(http.StatusOK, result)
	}

}

//// Получить все шаги для заданного по id сценария		(Сразу в PDF-е возвращать в JS???)
//func GetStepsFromScript(context *gin.Context)  {
//
//	helpers.SetLogFormat()
//	log.Infoln("Пришёл запрос в GetStepsFromScript")
//
//	// Данные из AJAX запроса
//	scriptId := context.PostForm("scriptId")
//	log.Infof("Данные из POST запроса AJAX: '%v' и '%v'", scriptId)
//
//	// Получить Шаги из БД только для заданных по ID Сценариев
//	scriptIdInt, _ := strconv.Atoi(scriptId)
//	scriptsIdList := append(make([]int, 1, 0), scriptIdInt)		// Слайс только из одного Id
//
//	stepsList, err := helpers.GetStepsListForSpecifiedScripts(scriptsIdList)
//	log.Debugf("%v - %v", stepsList, err)
//
//}














