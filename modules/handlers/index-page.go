package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../helpers"
	log "github.com/Sirupsen/logrus"
	"sort"
	"fmt"
)

var Version string = "1.6"

func ShowIndexPage(context *gin.Context)  {

	var testGroupList []string
	var err error

	testGroupList, err = helpers.GetTestGroupsList()

	if err != nil {
		log.Errorf("Ошибка при получении списка групп тестов из БД: %v", err)

		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": "Ошибка при получении списка групп тестов из БД",
				"message3": fmt.Sprintf("%s: ", err),
			},
		)

	} else {

		// Отсортировать список Групп
		sort.Strings(testGroupList)

		for idx, group := range testGroupList {
			log.Infof("%d.%s", idx, group)
		}

		// Обработка шаблона
		context.HTML(
			http.StatusOK,
			"index.html",
			gin.H{
				"title":         "SpecAdmin",
				"Version":       Version,
				"testGroupList": testGroupList,
			},
		)
	}
}
