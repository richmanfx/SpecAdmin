package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var Version string = "1.0"

func ShowIndexPage(context *gin.Context)  {

	// Обработка шаблона
	context.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title": "SpecAdmin",
			"Version": Version,
		},
	)
}
