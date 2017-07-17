package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"fmt"
	"net/http"
	log "github.com/Sirupsen/logrus"
	"../../models"
	"strconv"
	"os"
	"io"
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
	// Скриншот
	var screenShotFileName string
	screenShotFile, header, err := context.Request.FormFile("screen_shot")		// TODO: Переделать на простой FormFile !!!
	if err != nil {		// Если в форме не указан файл скриншота, то пустую строку - без скриншота
		log.Debugln("Не указан файл скриншота в функции 'AddStep' - передаём пустую строку (\"\").")
		screenShotFileName = ""
	} else {

		screenShotFileName = header.Filename
		log.Debugf("Загружается файл '%s'.", screenShotFileName)

		// Генерируем новое имя для изображения - скриншоты могут иметь одинаковое имя, храним уникальное
		screenShotFileName = helpers.GetUniqueFileName() + ".png"

		// Проверяем размер файла и если он превышает заданный размер
		// завершаем выполнение скрипта и выводим ошибку
		ScreenShotsSize, _ := strconv.Atoi(context.Request.Header.Get("Content-Length"))
		log.Debugf("Размер скриншота '%d' байт.", ScreenShotsSize)

		maxScreenShotsSize := 1000000 // Максимальный размер файла скриншота
		if ScreenShotsSize > maxScreenShotsSize {
			err = fmt.Errorf("Размер скриншота слишком большой - %d байт. Максимальный размер - %d байт.",
				ScreenShotsSize, maxScreenShotsSize)
		} else {
			log.Debugf("Данные из формы: newStepName='%v', stepSerialNumber='%v', stepsDescription='%v', " +
					   					"stepsExpectedResult='%v', stepsScript='%v', scriptsSuite='%v'",
				newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult, stepsScript, scriptsSuite)

			// Получить путь до хранилища скриншотов
			var screenShotsPath string
			screenShotsPath = helpers.GetScreenShotsPath()

			lastSymbolOfPath := screenShotsPath[len(screenShotsPath)-1:]
			log.Debugf("Последний символ в пути: '%s'", lastSymbolOfPath)
			var fullScreenShotsPath string
			if lastSymbolOfPath != string(os.PathSeparator) {
				fullScreenShotsPath = screenShotsPath + string(os.PathSeparator) + screenShotFileName
			} else {
				fullScreenShotsPath = screenShotsPath + screenShotFileName
			}

			log.Debugf("Полный путь к файлу скриншота: '%s'", fullScreenShotsPath)
			out, err := os.Create(fullScreenShotsPath)
			if err != nil {
				panic(err)
			}
			defer out.Close() // Файл закроется после работы с ним, даже при панике
			_, err = io.Copy(out, screenShotFile)
			if err != nil {
				panic(err)
			}
		}
	}

	// Группа Сюиты
	suite, err := helpers.GetSuite(scriptsSuite)
	suitesGroup := suite.Group

	if err == nil {
		// Добавить Шаг в БД
		err = helpers.AddTestStep(
			newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult, screenShotFileName, stepsScript, scriptsSuite)
	}

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
				"SuitesGroup":	suitesGroup,
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

	// Группа Сюиты
	suite, err := helpers.GetSuite(scriptsSuite)
	suitesGroup := suite.Group

	err = helpers.DelTestStep(deletedStep, stepsScript, scriptsSuite)

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
				"SuitesGroup":	suitesGroup,
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

	// Скриншот
	var screenShotFileName string
	screenShotFile, header, err := context.Request.FormFile("screen_shot")		// TODO: Переделать на простой FormFile !!!
	if err != nil {		// Если в форме не указан файл скриншота, то оставить прежний файл
		log.Infoln("Не указан файл скриншота.")
		screenShotFileName = ""
	} else {

		screenShotFileName = header.Filename
		log.Debugf("Загружается файл '%s'.", screenShotFileName)

		// Генерируем новое имя для изображения - скриншоты могут иметь одинаковое имя, храним уникальное
		screenShotFileName = helpers.GetUniqueFileName() + ".png"

		// Проверяем размер файла и если он превышает заданный размер
		// завершаем выполнение скрипта и выводим ошибку
		ScreenShotsSize, _ := strconv.Atoi(context.Request.Header.Get("Content-Length"))
		log.Debugf("Размер скриншота '%d' байт.", ScreenShotsSize)

		maxScreenShotsSize := 1000000 // Максимальный размер файла скриншота
		if ScreenShotsSize > maxScreenShotsSize {
			err = fmt.Errorf("Размер скриншота слишком большой - %d байт. Максимальный размер - %d байт.",
				ScreenShotsSize, maxScreenShotsSize)
		} else {
			log.Debugf("Данные из формы: stepsId='%v', stepsName='%v', stepsSerialNumber='%v', stepsDescription='%v', stepsExpectedResult='%v'",
				stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult)

			// Получить путь до хранилища скриншотов
			var screenShotsPath string
			screenShotsPath = helpers.GetScreenShotsPath()

			lastSymbolOfPath := screenShotsPath[len(screenShotsPath)-1:]
			log.Debugf("Последний символ в пути: '%s'", lastSymbolOfPath)
			var fullScreenShotsPath string
			if lastSymbolOfPath != string(os.PathSeparator) {
				fullScreenShotsPath = screenShotsPath + string(os.PathSeparator) + screenShotFileName
			} else {
				fullScreenShotsPath = screenShotsPath + screenShotFileName
			}

			log.Debugf("Полный путь к файлу скриншота: '%s'", fullScreenShotsPath)
			out, err := os.Create(fullScreenShotsPath)
			if err != nil {
				panic(err)
			}
			defer out.Close() // Файл закроется после работы с ним, даже при панике
			_, err = io.Copy(out, screenShotFile)
			if err != nil {
				panic(err)
			}
		}
	}


	// Обновить в БД
	err = helpers.UpdateTestStep(
		stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult, screenShotFileName)

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
	log.Debugln("Пришёл запрос в GetStepsOptions")

	// Данные из AJAX запроса
	stepsScriptsId, err := strconv.Atoi(context.PostForm("ScriptsId"))
	log.Debugf("Данные из POST запроса AJAX: stepsScriptsId='%d'", stepsScriptsId)
	if err != nil { panic(err) }

	// Данные о Шаге из БД
	stepsScriptName, scripsSuiteName, err := helpers.GetScriptAndSuiteByScriptId(stepsScriptsId)
	log.Debugf("Имя Сценария Шага: '%s'. Имя Сюиты Шага: '%s'.", stepsScriptName, scripsSuiteName)

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


// Удалить скриншот у Шага по заданному Id Шага (AJAX)
func DelScreenShotFromStep(context *gin.Context)  {
	helpers.SetLogFormat()
	log.Debugln("Пришёл запрос в DelScreenShotFromStep на удаление Скриншота.")

	// Данные из AJAX запроса
	stepsId, err := strconv.Atoi(context.PostForm("StepsId"))
	log.Debugf("Данные из POST запроса AJAX: stepsId='%d'", stepsId)
	if err != nil { panic(err) }

	// Удалить скриншот
	err = helpers.DeleteStepsScreenShot(stepsId)

	if err == nil {
		result := gin.H{"deleteStatus": "OK"}
		context.JSON(http.StatusOK, result)
	} else {
		log.Errorf("Ошибка удаления скриншота: %v", err)
		result := gin.H{"deleteStatus": err}
		context.JSON(http.StatusOK, result)
	}

}
