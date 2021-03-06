package handlers

import (
	"SpecAdmin/models"
	"SpecAdmin/modules/helpers"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Version = "10.4"

func ShowIndexPage(context *gin.Context) {
	var err error
	var statistic models.Statistic
	helpers.SetLogFormat()

	// Слайс из Групп
	groupList := make([]models.Group, 0, 40)

	// Сформировать список Групп на основе данных из БД
	err = helpers.GetGroupsList(&groupList)

	if err == nil {
		// Сбор статистики
		statistic, err = helpers.GetStatistic()
	}

	if err != nil {
		log.Errorf("Ошибка при получении списка групп тестов из БД: '%v'", err)
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":     "Ошибка",
				"message1":  "",
				"message2":  "Ошибка при получении списка групп тестов из БД",
				"message3":  fmt.Sprintf("%s: ", err),
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
			},
		)
	} else {
		context.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title":     "SpecAdmin",
				"Version":   Version,
				"UserLogin": helpers.UserLogin,
				"groupList": groupList,
				"statistic": statistic,
			},
		)
	}
}

func ShowSuitesIndex(context *gin.Context) {

	var err error
	var statistic models.Statistic
	var screenShotsPath string
	suitesList := make([]models.Suite, 0, 0) // Слайс из Сюит
	helpers.SetLogFormat()

	// Данные из формы
	groupName := context.PostForm("group_name")
	log.Debugf("Полученное из формы имя группы: '%s'", groupName)

	// Сформировать список Сюит заданной Группы из БД
	suitesList, err = helpers.GetSuitesListInGroup(groupName)

	if err == nil {
		// Сбор статистики
		statistic, err = helpers.GetStatistic()
	}

	if err == nil {
		// Путь к сриншотам
		screenShotsPath, err = helpers.GetScreenShotsPath()
	}

	if err != nil {
		log.Errorf("Ошибка при получении из БД списка Сюит для Группы тестов: %v", err)
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title":    "Ошибка",
				"message1": "",
				"message2": "Ошибка при получении из БД списка Сюит для Группы тестов",
				"message3": fmt.Sprintf("%s: ", err),
				"Version":  Version,
				"user":     helpers.UserLogin,
			},
		)
	} else {
		context.HTML(
			http.StatusOK,
			"suites-index.html",
			gin.H{
				"title":           "SpecAdmin",
				"Version":         Version,
				"UserLogin":       helpers.UserLogin,
				"groupName":       groupName,
				"suitesList":      suitesList,
				"statistic":       statistic,
				"screenShotsPath": screenShotsPath,
			},
		)
	}
}
