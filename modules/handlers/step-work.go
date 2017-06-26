package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"fmt"
	"net/http"
	_ "github.com/Sirupsen/logrus"
	_ "../../models"
	_ "strconv"
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





















