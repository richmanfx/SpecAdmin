package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	log "github.com/Sirupsen/logrus"
	"net/http"

)

func AuthRequired() gin.HandlerFunc {

	return func(context *gin.Context) {

		// Set example variable
		//context.Set("example", "12345")

		session := sessions.Default(context)	// Получить сессию контекста


		cookies := session.Get("Cookie")		// Получить из сессии Куки???

		// Получить отдельные куки
		//var splitCookie map[string]string
		//splitCookie = GetSplitCookie(cookies)

		//session.Delete("Cookie")	// Удалить Куки - на будущее

		// Если Кук нет, то устанавливаем их
		if cookies == nil{

			// Выставить Кук
			newCookie := "sessid=d232rn38jd1023e1nm13r25z;"
			session.Set("Cookie", newCookie)
			session.Save()

		} else {

			log.Printf("Куки от браузера: %v", cookies)
		}



		context.JSON(http.StatusOK, gin.H{"Cookie": cookies})



		// before request
		//context.Next()



		// after request
		//log.Printf("Здесь будем логиниться.")

		// access the status we are sending
		//status := context.Writer.Status()
		//log.Println(status)
	}
}


// Получить все Куки и вернуть их раздельно в МАПе
func GetSplitCookie(cookies string) map[string]string {

	var splitCookie map[string]string

	// Разделить по ";"


	// Разделить по "=" на ключ/значение


	return splitCookie
}