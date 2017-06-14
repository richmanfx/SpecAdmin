package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
)

// Добавить сюиту в базу
func AddSuite(context *gin.Context)  {


	newSuite := context.PostForm("suite")							// Сюита из формы
	suitesDescription := context.PostForm("suites_description")		// Описание Сюиты
	suitesGroup := context.PostForm("suites_group")					// Группа Сюиты
	err := helpers.AddTestSuite(newSuite, suitesDescription, suitesGroup)
	if err != nil {
		panic(err)
	}
}