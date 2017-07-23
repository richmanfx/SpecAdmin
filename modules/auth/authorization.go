package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"../handlers"
	"../helpers"
	"strings"
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
		log.Printf("splitCookie: '%v'", splitCookie)


		// Если Кук нет или если 'sessid' в БД не нашёлся, то на страницу авторизации
		sessidExist := helpers.SessionIdExistInBD(splitCookie["sessid"])

		if (len(splitCookie) == 0) || (sessidExist == false) {
			context.Abort()
			//Login(context)
			context.Redirect(http.StatusSeeOther, "/spec-admin/login")
		}



		//if cookies == nil{
		//
		//	// Выставить Кук
		//	newCookie := "sessid=хитрыйIDкуки;"
		//	session.Set("Cookie", newCookie)
		//	session.Save()
		//
		//} else {
		//
		//	log.Printf("Куки от браузера: %v", cookies)
		//	session.Delete("Cookie")	// Удалить Куки - на будущее
		//	log.Println("Куки удалены")
		//}
		//
		//
		//
		//context.JSON(http.StatusOK, gin.H{"Cookie": cookies})



		// before request
		//context.Next()



		// after request
		//log.Printf("Здесь будем логиниться.")

		// access the status we are sending
		//status := context.Writer.Status()
		//log.Println(status)
	}
}


// Получить все Куки и вернуть их раздельно и в МАПе
func GetSplitCookie(cookies string) map[string]string {

	splitCookie := make(map[string]string)
	var cookieList []string

	// Разделить по ";"
	cookieList = strings.Split(cookies, ";")
	cookieList = cookieList[0:len(cookieList)-1]
	log.Printf("cookieList: '%v'", cookieList)

	// Разделить по "=" на ключ/значение
	for _, cookie := range cookieList {
		log.Printf("Кука: %s", cookie)
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

	log.Infof("Пользователь: %s, Пароль: %s", userName, userPassword)
}


