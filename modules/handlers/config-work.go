package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"net/http"
)

// Отобразить страницу конфигурации для редактирования параметров
func EditConfig(context *gin.Context)  {
	var err error
	helpers.SetLogFormat()

	// Получить конфигурационные данные из БД
	config, err := helpers.GetConfig()
	if len(config) == 0 {
		err = fmt.Errorf("Нет конфигурационных данных в БД")
	}

	if err != nil {
		log.Errorf("Ошибка при получении конфигурационных данных из БД: %v", err)
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	"Ошибка при получении конфигурационных данных из БД",
				"message3": 	fmt.Sprintf("%s: ", err),
				"Version":	Version,

			},
		)
	} else {
		context.HTML(http.StatusOK,	"config.html",
			gin.H{
				"title":        "SpecAdmin",
				"config":	 	config,
				"Version":	Version,
			},
		)
	}
}


// Сохранить конфигурацию в БД
func SaveConfig(context *gin.Context)  {

	// Параметры из формы
	screenShotPath := context.PostForm("option_Путь к скриншотам")

	// Записать в БД
	err := helpers.SaveConfig(screenShotPath)

	if err != nil {
		// Ошибка при сохранении конфигурации
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintln("Ошибка при сохранении конфигурации"),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,

			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintln("Конфигурация успешно сохранена"),
				"message2": "",
				"message3": "",
				"Version":	Version,

			},
		)
	}
}


// Удалить из хранилища неиспользуемые в БД файлы скриншотов
func DelUnusedScreenShotsFile(context *gin.Context)  {

	// Получить список имён неиспользуемых файлов скриншотов
	unusedFileList, err := helpers.GetUnusedScreenShotsFileName()
	log.Infof("Неиспользуемые файлы скриншотов для удаления: '%v'", unusedFileList)
	countDeletedFile := len(unusedFileList)

	// Удалить в цикле файлы
	for _, deletedFile := range unusedFileList {
		err = helpers.DelScreenShotsFile(deletedFile)
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintln("Ошибка при удалении неиспользуемых файлов скриншотов"),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,

			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintln("Неиспользуемые файлы скриншотов успешно удалены"),
				"message2": "",
				"message3": fmt.Sprintf("Удалено %d файла(ов).", countDeletedFile),
				"Version":	Version,

			},
		)
	}
}


// Конфигурирование пользователей
func UsersConfig(context *gin.Context)  {

	var err error
	helpers.SetLogFormat()

	// Считать в БД пользователей и их пермишены
	users := helpers.GetUsekrs

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintln("Ошибка при получении данных о пользователях из БД"),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,

			},
		)
	} else {
		context.HTML(
			http.StatusOK,
			"users-config.html",
			gin.H{
				"title":   "SpecAdmin",
				"users":	users,
				"Version": Version,
			},
		)
	}
}
