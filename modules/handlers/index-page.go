package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../helpers"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"../../models"
)

var Version string = "2.7"

func ShowIndexPage(context *gin.Context)  {

	helpers.SetLogFormat()

	groupList := make([]models.Group, 0, 0)	// Слайс из Групп

	var err error

	// Сформировать список Групп на основе данных из БД
	groupList, err = helpers.GetGroupsList(groupList)

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
			},
		)
	}
}
