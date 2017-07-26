package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"../../models"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"net/http"
	//"strconv"
)

// Отобразить страницу конфигурации для редактирования параметров
func EditConfig(context *gin.Context)  {
	var err error
	helpers.SetLogFormat()

	// Получить конфигурационные данные из БД
	config, err := helpers.GetConfig()
	if len(config) == 0 {
		err = fmt.Errorf("Нет конфигурационных данных в БД")
	}

	if err != nil {
		log.Errorf("Ошибка при получении конфигурационных данных из БД: %v", err)
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	"Ошибка при получении конфигурационных данных из БД",
				"message3": 	fmt.Sprintf("%s: ", err),
				"Version":	Version,

			},
		)
	} else {
		context.HTML(http.StatusOK,	"config.html",
			gin.H{
				"title":        "SpecAdmin",
				"config":	 	config,
				"Version":	Version,
			},
		)
	}
}


// Сохранить конфигурацию в БД
func SaveConfig(context *gin.Context)  {

	// Параметры из формы
	screenShotPath := context.PostForm("option_Путь к скриншотам")

	// Записать в БД
	err := helpers.SaveConfig(screenShotPath)

	if err != nil {
		// Ошибка при сохранении конфигурации
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintln("Ошибка при сохранении конфигурации"),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,

			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintln("Конфигурация успешно сохранена"),
				"message2": "",
				"message3": "",
				"Version":	Version,

			},
		)
	}
}


// Удалить из хранилища неиспользуемые в БД файлы скриншотов
func DelUnusedScreenShotsFile(context *gin.Context)  {

	// Получить список имён неиспользуемых файлов скриншотов
	unusedFileList, err := helpers.GetUnusedScreenShotsFileName()
	log.Infof("Неиспользуемые файлы скриншотов для удаления: '%v'", unusedFileList)
	countDeletedFile := len(unusedFileList)

	// Удалить в цикле файлы
	for _, deletedFile := range unusedFileList {
		err = helpers.DelScreenShotsFile(deletedFile)
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintln("Ошибка при удалении неиспользуемых файлов скриншотов"),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,

			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintln("Неиспользуемые файлы скриншотов успешно удалены"),
				"message2": "",
				"message3": fmt.Sprintf("Удалено %d файла(ов).", countDeletedFile),
				"Version":	Version,

			},
		)
	}
}


// Конфигурирование пользователей
func UsersConfig(context *gin.Context)  {

	var err error
	helpers.SetLogFormat()

	// Считать из БД пользователей и их пермишены
	users, err := helpers.GetUsers()
	log.Infof("Пользователи из БД: '%v'", users)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintln("Ошибка при получении данных о пользователях из БД"),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,

			},
		)
	} else {
		context.HTML(
			http.StatusOK,
			"users-config.html",
			gin.H{
				"title":   "SpecAdmin",
				"users":	users,
				"Version": Version,
			},
		)
	}
}


// Создать нового пользователя
func CreateUser(context *gin.Context)  {

	var err error
	helpers.SetLogFormat()

	// Пользователь
	var user models.User

	// Данные из формы
	user.Login = context.PostForm("login")
	user.Password = context.PostForm("password")
	user.FullName = context.PostForm("full_name")

	if context.PostForm("create_permission") == "on" {
		user.Permissions.Create = true
	} else {
		user.Permissions.Create = false
	}

	if context.PostForm("edit_permission") == "on" {
		user.Permissions.Edit = true
	} else {
		user.Permissions.Edit = false
	}

	if context.PostForm("delete_permission") == "on" {
		user.Permissions.Delete = true
	} else {
		user.Permissions.Delete = false
	}

	if context.PostForm("config_permission") == "on" {
		user.Permissions.Config = true
	} else {
		user.Permissions.Config = false
	}

	if context.PostForm("users_permission") == "on" {
		user.Permissions.Users = true
	} else {
		user.Permissions.Users = false
	}

	log.Infof("user из формы = '%v'", user)

	// Создать пользователя в БД
	err = helpers.CreateUserInDb(user)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": "Ошибка при создании пользователя в БД.",
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Пользователь '%s' успешно добавлен в БД.", user.Login),
				"message2": "",
				"message3": "",
				"Version":	Version,
			},
		)
	}
}

























