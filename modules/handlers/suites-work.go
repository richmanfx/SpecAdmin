package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"../../models"
	"fmt"
	"net/http"
	log "github.com/Sirupsen/logrus"
)

// Добавить сюиту в базу
func AddSuite(context *gin.Context)  {

	helpers.SetLogFormat()

	newSuite := context.PostForm("suite")							// Сюита из формы
	suitesDescription := context.PostForm("suites_description")		// Описание Сюиты
	suitesSerialNumber := context.PostForm("suites_serial_number")	// Порядковый номер
	suitesGroup := context.PostForm("suites_group")					// Группа Сюиты
	err := helpers.AddTestSuite(newSuite, suitesDescription, suitesSerialNumber, suitesGroup)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при добавлении сюиты '%s' в группу '%s'.", newSuite, suitesGroup),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Сюита '%s' успешно добавлена в группу '%s'.", newSuite, suitesGroup),
				"message2": "",
				"message3": "",
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}

// Удалить сюиту из БД
func DelSuite(context *gin.Context)  {

	helpers.SetLogFormat()

	deletedSuite := context.PostForm("suite")							// Сюита из формы
	err := helpers.DelTestSuite(deletedSuite)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при удалении сюиты '%s'.", deletedSuite),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Сюита '%s' удалена успешно.", deletedSuite),
				"message2": "",	"message3": "",
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}

// Редактировать сюиту, данные получить из базы
func EditSuite(context *gin.Context)  {
	helpers.SetLogFormat()

	editedSuite := context.PostForm("suite")							// Сюита из формы
	log.Debugf("Редактируется сюита: %v", editedSuite)

	// Получить данные о сюите из БД
	var suite models.Suite
	var err error
	suite, err = helpers.GetSuite(editedSuite)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка получения данных о сюите '%s'.", editedSuite),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		// Вывести данные для редактирования
		log.Debugf("Сюита из БД: %v", suite)
		context.HTML(http.StatusOK, "edit-suite.html",
			gin.H{
				"title": 	"Редактирование сюиты",
				"suite": 	suite,
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}

// Обновить в БД сюиту после редактирования
func UpdateAfterEditSuite(context *gin.Context)  {
	helpers.SetLogFormat()

	//Данные из формы
	suitesName := context.PostForm("suite")
	suitesGroup := context.PostForm("suites_group")
	suitesSerialNumber := context.PostForm("suites_serial_number")
	suitesDescription := context.PostForm("suites_description")

	// Обновить в БД
	err := helpers.UpdateTestSuite(suitesName, suitesDescription, suitesSerialNumber, suitesGroup)
	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintf("Ошибка при обновлении сюиты '%s' в группе '%s'.", suitesName, suitesGroup),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Сюита '%s' успешно обновлена.", suitesName),
				"message2": "",
				"message3": "",
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}

}
