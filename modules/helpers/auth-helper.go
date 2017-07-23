package helpers

import (
	"time"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)


// Проверить наличие sessid в БД
func SessionIdExistInBD(sessidFromBrowser string) bool {

	var err error
	var sessidExist bool = false
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Если expire у сессии истекло, то сессию удалить из БД и вернуть false
	var expires time.Time

	requestResult := db.QueryRow("SELECT expires FROM sessions WHERE session_id=?", sessidFromBrowser)
	log.Infof("requestResult: %v", requestResult)
	err = requestResult.Scan(&expires)

	if err == nil {
		log.Infof("'expires' из БД для sessid='%s': %v", sessidFromBrowser, expires)

		// Если expire позднее текущего времени, то сессия существует и валидна
		if expires.After(time.Now()) {
			log.Infoln("Сессия не истекла")
			sessidExist = true
		} else {
			// Если сессия истекла, то удаляем её из БД
			result, err := db.Exec("DELETE FROM sessions WHERE session_id=?", sessidFromBrowser)
			if err == nil {
				var affected int64
				affected, err = result.RowsAffected()
				if err == nil {
					if affected == 0 {
						err = fmt.Errorf("Ошибка при удалениии сессии '%s'.", sessidFromBrowser)
						log.Debugf("Сессия '%s' НЕ удалена.", sessidFromBrowser)
					}
					log.Debugf("Удалено '%d' строк в таблице 'sessions'.", affected)
				}
			}
			sessidExist = false
		}
	} else {
		// Такой сессии в БД нет
		log.Infof("Сессии '%s' в таблице 'sessions' нет: %v", sessidFromBrowser, err)
		sessidExist = false
	}

	db.Close()

	return sessidExist
}


// Сохранить в БД Сессию
func SaveSessionInDB(sessid string, expires time.Time, userName string) error {

	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }


	// Внести в БД
	result, err := db.Exec("INSERT INTO sessions (session_id,expires,user) VALUE (?,?,?)", sessid, expires, userName)
	//result, err := db.Exec("INSERT INTO sessions (session_id, expires) VALUE (?,?)", sessid, expires)

	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {
			panic(err)
		}
		log.Debugf("Вставлено %d строк в таблицу 'sessions'.", affected)
	}

	db.Close()
	return err
}