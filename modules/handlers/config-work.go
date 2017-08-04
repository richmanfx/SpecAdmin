package handlers

import (
	"github.com/gin-gonic/gin"
	"../helpers"
	"../../models"
	log "github.com/Sirupsen/logrus"
	"fmt"
	"net/http"
)

// Отобразить страницу конфигурации для редактирования параметров
func EditConfig(context *gin.Context)  {
	var err error
	var config []models.Option
	helpers.SetLogFormat()

	// Проверить пермишен пользователя для работы с конфигурацией
	log.Debugf("Проверка пермишена для пользователя '%s'", helpers.UserLogin)
	err = helpers.CheckOneUserPermission(helpers.UserLogin, "config_permission")

	if err == nil {
		// Получить конфигурационные данные из БД
		config, err = helpers.GetConfig()
		if len(config) == 0 {
			err = fmt.Errorf("Нет конфигурационных данных в БД")
		}
	}

	if err != nil {
		log.Errorf("Ошибка при получении конфигурационных данных из БД: %v", err)
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	"Ошибка при получении конфигурационных данных из БД",
				"message3": 	fmt.Sprintf("%s: ", err),
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK,	"config.html",
			gin.H{
				"title":        "SpecAdmin",
				"config":	 	config,
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
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
				"UserLogin":	helpers.UserLogin,
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
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}


// Удалить из хранилища неиспользуемые в БД файлы скриншотов
func DelUnusedScreenShotsFile(context *gin.Context)  {

	var countDeletedFile int

	// Проверить пермишен пользователя для удалений
	log.Debugf("Проверка пермишена для пользователя '%s'", helpers.UserLogin)
	err := helpers.CheckOneUserPermission(helpers.UserLogin, "delete_permission")

	if err == nil {

		// Получить список имён неиспользуемых файлов скриншотов
		var unusedFileList []string
		unusedFileList, err = helpers.GetUnusedScreenShotsFileName()
		log.Debugf("Неиспользуемые файлы скриншотов для удаления: '%v'", unusedFileList)
		countDeletedFile = len(unusedFileList)

		// Удалить в цикле файлы
		for _, deletedFile := range unusedFileList {
			err = helpers.DelScreenShotsFile(deletedFile)
		}
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": fmt.Sprintln("Ошибка при удалении неиспользуемых файлов скриншотов"),
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
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
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}


// Конфигурирование пользователей
func UsersConfig(context *gin.Context)  {

	var err error
	var users []models.User
	helpers.SetLogFormat()

	// Проверить пермишен пользователя для работы с пользователями
	log.Debugf("Проверка пермишена для пользователя '%s'", helpers.UserLogin)
	err = helpers.CheckOneUserPermission(helpers.UserLogin, "users_permission")

	if err == nil {
		// Считать из БД пользователей и их пермишены
		users, err = helpers.GetUsers()
		log.Debugf("Пользователи из БД: '%v'", users)
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	fmt.Sprintln("Ошибка при получении данных о пользователях из БД"),
				"message3": 	fmt.Sprintf("%s: ", err),
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
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
				"UserLogin":	helpers.UserLogin,
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
	openPassword := context.PostForm("password")
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

	log.Debugf("user из формы создания = '%v'", user)


	// Получить Соль и Хеш пароля
	user.Salt = helpers.CreateSalt()
	user.Password = helpers.CreateHash(openPassword, user.Salt)

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
				"UserLogin":	helpers.UserLogin,
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
				"UserLogin":	helpers.UserLogin,
			},
		)
	}
}


// Сохранить пользователя после редактирования
func SaveUser(context *gin.Context)  {

	log.Debugln("Мы в 'SaveUser'!")

	var err error
	helpers.SetLogFormat()

	// Пользователь
	var savedUser models.User

	// Данные из формы
	savedUser.Login = context.PostForm("login")
	savedUser.FullName = context.PostForm("full_name")

	if context.PostForm("create_permission") == "on" {
		savedUser.Permissions.Create = true
	} else {
		savedUser.Permissions.Create = false
	}

	if context.PostForm("edit_permission") == "on" {
		savedUser.Permissions.Edit = true
	} else {
		savedUser.Permissions.Edit = false
	}

	if context.PostForm("delete_permission") == "on" {
		savedUser.Permissions.Delete = true
	} else {
		savedUser.Permissions.Delete = false
	}

	if context.PostForm("config_permission") == "on" {
		savedUser.Permissions.Config = true
	} else {
		savedUser.Permissions.Config = false
	}

	if context.PostForm("users_permission") == "on" {
		savedUser.Permissions.Users = true
	} else {
		savedUser.Permissions.Users = false
	}

	log.Debugf("savedUser из формы редактирования = '%v'", savedUser)

	// Сохранить пользователя в БД
	err = helpers.SaveUserInDb(savedUser)

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": 	"",
				"message2": 	"Ошибка при сохранении пользователя в БД.",
				"message3": 	fmt.Sprintf("%s: ", err),
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Информация",
				"message1": 	fmt.Sprintf("Пользователь '%s' успешно сохранён в БД.", savedUser.Login),
				"message2": 	"",
				"message3": 	"",
				"Version":		Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}

}

// Удалить пользователя
func DeleteUser(context *gin.Context)  {

	var err error
	var deletedUser models.User // Пользователь
	helpers.SetLogFormat()

	// Проверить пермишен пользователя для удалений
	log.Debugf("Проверка пермишена для пользователя '%s'", helpers.UserLogin)
	err = helpers.CheckOneUserPermission(helpers.UserLogin, "delete_permission")

	if err == nil {



		// Данные из формы
		deletedUser.Login = context.PostForm("login")
		deletedUser.FullName = context.PostForm("full_name")

		log.Debugf("deletedUser из формы удаления = '%v'", deletedUser)

		// Удалить из БД
		err = helpers.DeleteUserInDb(deletedUser)
	}

	if err != nil {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Ошибка",
				"message1": "",
				"message2": "Ошибка при удалении пользователя из БД.",
				"message3": fmt.Sprintf("%s: ", err),
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintf("Пользователь '%s' успешно удалён из БД.", deletedUser.Login),
				"message2": "",
				"message3": "",
				"Version":	Version,
				"UserLogin":	helpers.UserLogin,
			},
		)
	}

}


























