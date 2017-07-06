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
			},
		)
	} else {
		context.HTML(http.StatusOK,	"config.html",
			gin.H{
				"title":        "SpecAdmin",
				"Version":      Version,
				"config":	 	config,
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
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintln("Конфигурация успешно сохранена"),
				"message2": "",
				"message3": "",
			},
		)
	}
}
