package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
)

// Добавить сюиту в базу
func AddSuite(context *gin.Context)  {


	newSuite := context.PostForm("suite")		// Сюита из формы
	err := helpers.AddTestSuite(newSuite)
	if err != nil {
		panic(err)
	}
}