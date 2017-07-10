package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../helpers"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"../../models"
)

var Version string = "5.2"

func ShowIndexPage(context *gin.Context)  {
	var err error
	helpers.SetLogFormat()

	// Слайс из Групп
	groupList := make([]models.Group, 0, 0)

	// Сформировать список Групп на основе данных из БД
	groupList, err = helpers.GetGroupsList(groupList)

	// Сбор статистики
	statistic, err := helpers.GetStatistic()

	if err != nil {
		log.Errorf("Ошибка при получении списка групп тестов из БД: %v", err)
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	"Ошибка при получении списка групп тестов из БД",
				"message3": 	fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title":        "SpecAdmin",
				"Version":      Version,
				"groupList": 	groupList,
				"statistic":	statistic,
			},
		)
	}
}

func ShowSuitesIndex(context *gin.Context)  {

	var err error
	helpers.SetLogFormat()

	// Данные из формы
	groupName := context.PostForm("group_name")
	log.Debugf("Полученное из формы имя группы: '%s'", groupName)

	// Слайс из Сюит
	suitesList := make([]models.Suite, 0, 0)

	// Сформировать список Сюит заданной Группы из БД
	suitesList, err = helpers.GetSuitesListInGroup(groupName)

	// Сбор статистики
	statistic, err := helpers.GetStatistic()

	// Путь к сриншотам
	screenShotsPath := helpers.GetScreenShotsPath()

	if err != nil {
		log.Errorf("Ошибка при получении из БД списка Сюит для Группы тестов.: %v", err)
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	"Ошибка при получении из БД списка Сюит для Группы тестов",
				"message3": 	fmt.Sprintf("%s: ", err),
			},
		)
	} else {
		context.HTML(
			http.StatusOK,
			"suites-index.html",
			gin.H{
				"title":        "SpecAdmin",
				"Version":      Version,
				"groupName":	groupName,
				"suitesList": 	suitesList,
				"statistic":	statistic,
				"screenShotsPath":	screenShotsPath,
			},
		)
	}
}
