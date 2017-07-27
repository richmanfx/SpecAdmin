package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"../handlers"
	"../helpers"
	"strings"
	"fmt"
	"time"
)

// Мидлеваря - проверяет авторизацию
func AuthRequired() gin.HandlerFunc {

	return func(context *gin.Context) {

		// Set example variable
		//context.Set("example", "12345")

		session := sessions.Default(context)	// Получить сессию контекста

		cookies := session.Get("Cookie")		// Получить из сессии все Куки


		// Получить отдельные куки
		var splitCookie map[string]string
		splitCookie = GetSplitCookie(cookies.(string))
		log.Infof("Разделённые Куки из браузера: '%v'", splitCookie)


		// Если Кук нет или если 'sessid' в БД не нашёлся, то на страницу авторизации
		sessidValue := splitCookie["sessid"]
		sessidExist := helpers.SessionIdExistInBD(sessidValue)

		if len(splitCookie) == 0 {
			log.Infoln("Кук из браузера нет.")
			context.Abort()
			context.Redirect(http.StatusSeeOther, "/spec-admin/login")
		} else if sessidExist == false {
			log.Infoln("'sessid' из браузера не обнаружена в БД.")
			context.Abort()
			context.Redirect(http.StatusSeeOther, "/spec-admin/login")
		}
	}
}


// Получить все Куки и вернуть их раздельно и в МАПе
func GetSplitCookie(cookies string) map[string]string {

	splitCookie := make(map[string]string)
	var cookieList []string

	// Разделить по ";"
	cookieList = strings.Split(cookies, ";")
	cookieList = cookieList[0:len(cookieList)-1]
	log.Infof("Список Кук из браузера: '%v'", cookieList)

	// Разделить по "=" на ключ/значение
	for _, cookie := range cookieList {
		log.Debugf("Кука: %s", cookie)
		splitCookie[strings.Split(cookie, "=")[0]] = strings.Split(cookie, "=")[1]
	}

	return splitCookie
}

// Страница авторизации
func Login(context *gin.Context)  {
	context.HTML(
		http.StatusOK,
		"login.html",
		gin.H{
			"title":        "SpecAdmin",
			"Version":      handlers.Version,
		},
	)
}

// Авторизация пользователя
func Authorization(context *gin.Context)  {
	userName := context.PostForm("user_name")
	userPassword := context.PostForm("user_password")
	log.Debugf("Пользователь: %s, Пароль: %s", userName, userPassword)

	// Пока без базы проверяем
	validUserName := "user"
	validUserPassword := "Qwerty123"
	if (userName == validUserName) && (userPassword == validUserPassword) {

		// Изменить сессию
		session := sessions.Default(context)

		// Сгенерировать sessid
		sessid := helpers.GetUniqueFileName()		// TODO: временно, сделать отдельный генератор
		newCookie := fmt.Sprintf("sessid=%s;", sessid)
		log.Infof("Сгенерирована новая Кука: '%s'", newCookie)


		//	session.Delete("Cookie")	// Удалить Куки - на будущее
		//	log.Println("Куки удалены")
		//context.JSON(http.StatusOK, gin.H{"Cookie": cookies})

		// Сохранить сессию в БД
		var expire time.Time = time.Now().Add(12 * time.Hour)	// Кука устаревает через 12 часов
		err := helpers.SaveSessionInDB(sessid, expire, validUserName)
		if err != nil {
			log.Errorf("Ошибка сохранения сессии в БД: %v", err)
			context.HTML(http.StatusOK, "message.html",
				gin.H{
					"title": 		"Ошибка",
					"message1": 	"",
					"message2": 	"Ошибка сохранения сессии в БД.",
					"message3": 	err,
					"Version":		handlers.Version,
				},
			)
		}

		// Выставить в браузере Куки
		session.Set("Cookie", newCookie)
		session.Save()
		log.Infof("Новая Кука '%s' отправлена в браузер", newCookie)

		// Направить на индексную страницу
		context.Abort()
		context.Redirect(http.StatusSeeOther, "/spec-admin")

	} else {
		log.Errorln("Ошибка авторизации - неверный логин/пароль.")
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	"Ошибка авторизации - неверный логин/пароль.",
				"message3": 	"",
				"Version":		handlers.Version,
			},
		)
	}


}


// Разлогирование
func Logout(context *gin.Context)  {

	context.Abort()
	context.Redirect(http.StatusSeeOther, "/spec-admin/login")
}


// Смена пароля
func ChangePassword(context *gin.Context)  {

	// Данные из формы
	userName := context.PostForm("login")
	oldPassword := context.PostForm("old_password")
	newPassword := context.PostForm("new_password")
	log.Infof("Данные из формы смены пароля: '%s', '%s', '%s'", userName, oldPassword, newPassword)


	// Проверить валидность старого пароля
	err := helpers.ValidatePassword(userName, oldPassword)
	if err != nil {
		log.Errorln("Указан неверный старый пароль.")
	} else {
		// Записать в БД новый пароль
		err = helpers.SavePassword(userName, newPassword)
	}

	if err != nil {
		log.Errorf("Ошибка изменения пароля: %v", err)
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": 		"Ошибка",
				"message1": 	"",
				"message2": 	"Ошибка изменения пароля.",
				"message3": 	err,
				"Version":		handlers.Version,
			},
		)
	} else {
		context.HTML(http.StatusOK, "message.html",
			gin.H{
				"title": "Информация",
				"message1": fmt.Sprintln("Пароль успешно изменён"),
				"message2": "",
				"message3": "",
				"Version":	handlers.Version,
			},
		)
	}

}

