package handlers

import (
	"SpecAdmin/models"
	"SpecAdmin/modules/helpers"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

var stepBuffer int

// Добавить новый Шаг
func AddStep(context *gin.Context) {

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
	screenShotFile, header, err := context.Request.FormFile("screen_shot") // TODO: Переделать на простой FormFile !!!
	if err != nil {                                                        // Если в форме не указан файл скриншота, то пустую строку - без скриншота
		log.Debugln("Не указан файл скриншота в функции 'AddStep' - передаём пустую строку (\"\").")
		screenShotFileName = ""
		err = nil
	} else {

		screenShotFileName = header.Filename
		log.Debugf("Загружается файл '%s'.", screenShotFileName)

		// Генерируем новое имя для изображения - скриншоты могут иметь одинаковое имя, храним уникальное
		screenShotFileName = helpers.GetUnique32SymbolsString() + ".png"

		// Проверяем размер файла и если он превышает заданный размер
		// завершаем выполнение скрипта и выводим ошибку
		ScreenShotsSize, _ := strconv.Atoi(context.Request.Header.Get("Content-Length"))
		log.Debugf("Размер скриншота '%d' байт.", ScreenShotsSize)

		maxScreenShotsSize := 1000000 // Максимальный размер файла скриншота
		if ScreenShotsSize > maxScreenShotsSize {
			err = errors.New(fmt.Sprintf("Размер скриншота слишком большой - %d байт. Максимальный размер - %d байт.",
				ScreenShotsSize, maxScreenShotsSize))
		} else {
			log.Debugf("Данные из формы: newStepName='%v', stepSerialNumber='%v', stepsDescription='%v', "+
				"stepsExpectedResult='%v', stepsScript='%v', scriptsSuite='%v'",
				newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult, stepsScript, scriptsSuite)

			// Получить путь до хранилища скриншотов
			var screenShotsPath string
			screenShotsPath, err = helpers.GetScreenShotsPath()

			if err == nil {

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
				if err == nil {
					defer out.Close() // Файл закроется после работы с ним, даже при панике
					_, err = io.Copy(out, screenShotFile)
				}
			}
		}
	}

	var suitesGroup string
	var suite models.Suite
	if err == nil {
		// Группа Сюиты
		suite, err = helpers.GetSuite(scriptsSuite)
		suitesGroup = suite.Group
	}

	if err == nil {
		// Добавить Шаг в БД
		err = helpers.AddTestStep(
			newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult, screenShotFileName, stepsScript, scriptsSuite)
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":    "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при добавлении шага '%s' в сценарий '%s' сюиты '%s'.",
					newStepName, stepsScript, scriptsSuite),
				"message3":  fmt.Sprintf("%s: ", err),
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Шаг '%s' успешно добавлен в сценарий '%s' сюиты '%s'.",
					newStepName, stepsScript, scriptsSuite),
				"message2":    "",
				"message3":    "",
				"SuitesGroup": suitesGroup,
				"Version":     Version,
				"UserLogin":   helpers.UserLogin,
			},
		)
	}
}

// Удалить Шаг
func DelStep(context *gin.Context) {
	helpers.SetLogFormat()

	// Данные из HTML-формы
	deletedStep := context.PostForm("step")
	stepsScript := context.PostForm("script")
	scriptsSuite := context.PostForm("suite")

	// Группа Сюиты
	var suitesGroup string
	suite, err := helpers.GetSuite(scriptsSuite)
	if err == nil {
		suitesGroup = suite.Group
		err = helpers.DelTestStep(deletedStep, stepsScript, scriptsSuite)
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":     "Ошибка",
				"message1":  "",
				"message2":  fmt.Sprintf("Ошибка при удалении шага '%s'.", deletedStep),
				"message3":  fmt.Sprintf("%s: ", err),
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":    "Информация",
				"message1": fmt.Sprintf("Шаг '%s' успешно удалён.", deletedStep),
				"message2": "", "message3": "",
				"SuitesGroup": suitesGroup,
				"Version":     Version,
				"UserLogin":   helpers.UserLogin,
			},
		)
	}
}

// Редактировать Шаг
func EditStep(context *gin.Context) {
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
				"title":     "Ошибка",
				"message1":  "",
				"message2":  fmt.Sprintf("Ошибка получения данных о шаге '%s'.", editedStepName),
				"message3":  fmt.Sprintf("%s: ", err),
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
			},
		)
	} else {
		// Вывести данные для редактирования
		context.HTML(http.StatusOK, "edit-step.html",
			gin.H{
				"title":     "Редактирование шага",
				"step":      step,
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
			},
		)
	}
}

// Обновить в БД Шаг после редактирования
func UpdateAfterEditStep(context *gin.Context) {
	helpers.SetLogFormat()
	var stepsId, stepsSerialNumber int
	var stepsName, stepsDescription, stepsExpectedResult string
	var err error
	var screenShotFile multipart.File
	var header *multipart.FileHeader
	var screenShotFileName string

	//Данные из формы
	stepsId, err = strconv.Atoi(context.PostForm("hidden_id"))
	if err == nil {
		stepsName = context.PostForm("step")
		stepsSerialNumber, err = strconv.Atoi(context.PostForm("steps_serial_number"))
		if err == nil {
			stepsDescription = context.PostForm("steps_description")
			stepsExpectedResult = context.PostForm("steps_expected_result")

			// Скриншот
			screenShotFile, header, err = context.Request.FormFile("screen_shot") // TODO: Переделать на простой FormFile !!!
			if err != nil {                                                       // Если в форме не указан файл скриншота, то оставить прежний файл
				log.Debugln("Не указан файл скриншота.")
				screenShotFileName = ""
				err = nil
			} else {

				screenShotFileName = header.Filename
				log.Debugf("Загружается файл '%s'.", screenShotFileName)

				// Генерируем новое имя для изображения - скриншоты могут иметь одинаковое имя, храним уникальное
				screenShotFileName = helpers.GetUnique32SymbolsString() + ".png"

				// Проверяем размер файла и если он превышает заданный размер
				// завершаем выполнение скрипта и выводим ошибку
				ScreenShotsSize, _ := strconv.Atoi(context.Request.Header.Get("Content-Length"))
				log.Debugf("Размер скриншота '%d' байт.", ScreenShotsSize)

				maxScreenShotsSize := 1000000 // Максимальный размер файла скриншота
				if ScreenShotsSize > maxScreenShotsSize {
					err = errors.New(fmt.Sprintf("Размер скриншота слишком большой - %d байт. Максимальный размер - %d байт.",
						ScreenShotsSize, maxScreenShotsSize))
				} else {
					log.Debugf("Данные из формы: stepsId='%v', stepsName='%v', stepsSerialNumber='%v', stepsDescription='%v', stepsExpectedResult='%v'",
						stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult)

					// Получить путь до хранилища скриншотов
					var screenShotsPath string
					screenShotsPath, err = helpers.GetScreenShotsPath()

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
					if err == nil {
						defer out.Close() // Файл закроется после работы с ним, даже при панике
						_, err = io.Copy(out, screenShotFile)
					}
				}
			}
		}
	}

	ScriptId, err := helpers.GetScriptsIdByStepsId(stepsId)
	var suitesGroup string
	if err == nil {
		_, scriptsSuite, err := helpers.GetScriptAndSuiteByScriptId(ScriptId)
		if err == nil {

			var suite models.Suite
			// Группа Сюиты
			suite, err = helpers.GetSuite(scriptsSuite)
			suitesGroup = suite.Group

			if err == nil {
				// Обновить в БД
				err = helpers.UpdateTestStep(
					stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult, screenShotFileName)
			}
		}
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":     "Ошибка",
				"message1":  "",
				"message2":  fmt.Sprintf("Ошибка при обновлении Шага '%s'", stepsName),
				"message3":  fmt.Sprintf("%s: ", err),
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":       "Информация",
				"message1":    fmt.Sprintf("Шаг '%s' успешно обновлён.", stepsName),
				"message2":    "",
				"message3":    "",
				"SuitesGroup": suitesGroup,
				"Version":     Version,
				"UserLogin":   helpers.UserLogin,
			},
		)
	}
}

// Вернуть параметры Шага (AJAX)
func GetStepsOptions(context *gin.Context) {
	helpers.SetLogFormat()
	log.Debugln("Пришёл запрос в GetStepsOptions")
	var stepsScriptName, scripsSuiteName string

	// Данные из AJAX запроса
	stepsScriptsId, err := strconv.Atoi(context.PostForm("ScriptsId"))
	log.Debugf("Данные из POST запроса AJAX: stepsScriptsId='%d'", stepsScriptsId)
	if err == nil {
		// Данные о Шаге из БД
		stepsScriptName, scripsSuiteName, err = helpers.GetScriptAndSuiteByScriptId(stepsScriptsId)
		log.Debugf("Имя Сценария Шага: '%s'. Имя Сюиты Шага: '%s'.", stepsScriptName, scripsSuiteName)
	}

	if err == nil {
		result := gin.H{"stepsScriptName": stepsScriptName, "scripsSuiteName": scripsSuiteName}
		context.JSON(http.StatusOK, result)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":     "Ошибка",
				"message1":  "",
				"message2":  fmt.Sprintln("Ошибка в функции 'GetStepsOptions' при ответе на AJAX запрос"),
				"message3":  fmt.Sprintf("%s: ", err),
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
			},
		)
	}
}

// Удалить скриншот у Шага по заданному Id Шага (AJAX)
func DelScreenShotFromStep(context *gin.Context) {
	helpers.SetLogFormat()
	log.Debugln("Пришёл запрос в DelScreenShotFromStep на удаление Скриншота.")

	// Данные из AJAX запроса
	stepsId, err := strconv.Atoi(context.PostForm("StepsId"))
	log.Debugf("Данные из POST запроса AJAX: stepsId='%d'", stepsId)
	if err == nil {
		// Удалить скриншот
		err = helpers.DeleteStepsScreenShot(stepsId)
	}

	if err == nil {
		result := gin.H{"deleteStatus": "OK"}
		context.JSON(http.StatusOK, result)
	} else {
		log.Errorf("Ошибка удаления скриншота: %v", err)
		result := gin.H{"deleteStatus": err}
		context.JSON(http.StatusOK, result)
	}
}

// Поместить Шаг в буфер обмена
func CopyStepInClipboard(context *gin.Context) {
	helpers.SetLogFormat()
	log.Infoln("Пришёл запрос в CopyStepInClipboard для помещения Шага в буфер обмена.")

	// Данные из AJAX запроса
	stepId, err := strconv.Atoi(context.PostForm("StepId"))
	log.Infof("Данные из POST запроса AJAX: stepId='%d'", stepId)

	// Идентификатор Шага поместить в буфер
	stepBuffer = stepId

	if err == nil {
		result := gin.H{"Status": "OK"}
		context.JSON(http.StatusOK, result)
	} else {
		log.Errorf("Ошибка помещения Шага в буфер обмена: %v", err)
		result := gin.H{"Status": err}
		context.JSON(http.StatusOK, result)
	}

}

// Получить имя, описание и ожидаемый результат Шага
func GetStepFromBuffer(context *gin.Context) {

	var stepsName, stepsDescription, stepsExpectedResult string
	helpers.SetLogFormat()
	log.Debugln("Пришёл запрос в GetStepFromBuffer для получения информации о Шаге из буфера обмена.")

	// Данные из AJAX запроса
	stepId, err := strconv.Atoi(context.PostForm("StepsId"))
	log.Debugf("Данные из POST запроса AJAX: StepsId='%d'", stepId)

	// Получить данные о Шаге
	if err == nil {
		stepsName, stepsDescription, stepsExpectedResult, err = helpers.GetStepsData(stepBuffer)
	}

	if err == nil {
		result := gin.H{
			"stepsName":           stepsName,
			"stepsDescription":    stepsDescription,
			"stepsExpectedResult": stepsExpectedResult,
		}
		context.JSON(http.StatusOK, result)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":     "Ошибка",
				"message1":  "",
				"message2":  fmt.Sprintln("Ошибка в функции 'GetStepFromBuffer' при ответе на AJAX запрос"),
				"message3":  fmt.Sprintf("%s: ", err),
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
			},
		)
	}
}
