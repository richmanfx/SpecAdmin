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

// Добавить новый Шаг
func AddStep(context *gin.Context)  {

	helpers.SetLogFormat()

	// Данные из формы
	newStepName := context.PostForm("step")
	stepSerialNumber := context.PostForm("step_serial_number")
	stepsDescription := context.PostForm("steps_description")
	stepsExpectedResult := context.PostForm("steps_expected_result")
	stepsScript := context.PostForm("steps_script")
	scriptsSuite := context.PostForm("scripts_suite")

	err := helpers.AddTestStep(newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult, stepsScript, scriptsSuite)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при добавлении шага '%s' в сценарий '%s' сюиты '%s'.",
					newStepName, stepsScript, scriptsSuite),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Шаг '%s' успешно добавлен в сценарий '%s' сюиты '%s'.",
					newStepName, stepsScript, scriptsSuite),
				"message2": "",
				"message3": "",
			},
		)
	}
}


// Удалить Шаг
func DelStep(context *gin.Context)  {
	helpers.SetLogFormat()

	// Данные из HTML-формы
	deletedStep := context.PostForm("step")
	stepsScript := context.PostForm("script")
	scriptsSuite := context.PostForm("suite")

	err := helpers.DelTestStep(deletedStep, stepsScript, scriptsSuite)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при удалении шага '%s'.", deletedStep),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Шаг '%s' успешно удалён.", deletedStep),
				"message2": "",	"message3": "",
			},
		)
	}
}


// Редактировать Шаг
func EditStep(context *gin.Context)  {
	helpers.SetLogFormat()

	// Данные из формы
	editedStepName := context.PostForm("step")
	stepsScript := context.PostForm("steps_script")
	scriptsSuite := context.PostForm("scripts_suite")

	log.Debugf("Редактируется Шаг '%s' сценария '%s' в сюите '%s'.", editedStepName, stepsScript, scriptsSuite)

	// Получить данные о шаге из БД
	var err error
	var step models.Step
	step, err = helpers.GetStep(editedStepName, stepsScript, scriptsSuite)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка получения данных о шаге '%s'.", editedStepName),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		// Вывести данные для редактирования
		context.HTML(http.StatusOK, "edit-step.html",
			gin.H{
				"title": 	"Редактирование шага",
				"Version":	Version,
				"step": 	step,
			},
		)
	}
}


// Обновить в БД Шаг после редактирования
func UpdateAfterEditStep(context *gin.Context)  {
	helpers.SetLogFormat()

	//Данные из формы
	stepsId, err := strconv.Atoi(context.PostForm("hidden_id"))
	if err != nil { panic(err) }
	stepsName := context.PostForm("step")
	stepsSerialNumber, err := strconv.Atoi(context.PostForm("steps_serial_number"))
	if err != nil { panic(err) }
	stepsDescription := context.PostForm("steps_description")
	stepsExpectedResult := context.PostForm("steps_expected_result")
	log.Debugf("Данные из формы: stepsId='%v', stepsName='%v', stepsSerialNumber='%v', stepsDescription='%v', stepsExpectedResult='%v'",
		stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult)

	// Обновить в БД
	err = helpers.UpdateTestStep(stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при обновлении Шага '%s'", stepsName),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Шаг '%s' успешно обновлён.", stepsName),
				"message2": "",
				"message3": "",
			},
		)
	}
}

// Вернуть параметры Шага (AJAX)
func GetStepsOptions(context *gin.Context)  {
	helpers.SetLogFormat()
	log.Infoln("Пришёл запрос в GetStepsOptions")

	// Данные из формы
	stepsScriptsId, err := strconv.Atoi(context.PostForm("ScriptsId"))
	log.Infof("Данные из POST запроса AJAX: '%d'", stepsScriptsId)
	if err != nil { panic(err) }

	// Данные о Шаге из БД
	stepsScriptName, scripsSuiteName, err := helpers.GetScriptAndSuiteByScriptId(stepsScriptsId)

	if err == nil {
		result := gin.H{"stepsScriptName": stepsScriptName, "scripsSuiteName": scripsSuiteName}
		context.JSON(http.StatusOK, result)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintln("Ошибка в функции 'GetStepsOptions' при ответе на AJAX запрос"),
				"message3": fmt.Sprintf("%s: ", err),
			},
		)
	}


}
