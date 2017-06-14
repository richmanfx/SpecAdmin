package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../helpers"
	log "github.com/Sirupsen/logrus"
	"sort"
)

var Version string = "1.5"

func ShowIndexPage(context *gin.Context)  {

	var testGroupList []string = helpers.GetTestGroupsList()

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
			"title": "SpecAdmin",
			"Version": Version,
			"testGroupList": testGroupList,
		},
	)
}
