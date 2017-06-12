package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"../helpers"
	"fmt"
)

var Version string = "1.3"

func ShowIndexPage(context *gin.Context)  {


	testGroupList := helpers.GetTestGroupsList()

	for idx, group := range testGroupList {
		fmt.Printf("%d. %s\n", idx, group)
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
