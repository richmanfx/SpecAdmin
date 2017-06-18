package helpers

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Соединение с БД и проверка соединения
func dbConnect() error {

	var err error

	SetLogFormat()

	// Соединение с БД
	log.Debugf("Соединение с БД")
	db, err = sql.Open("mysql", "specuser:Ghashiz7@tcp(localhost:3306)/specadmin?charset=utf8")
	if err != nil {
		log.Errorf("Ошибка соединение с БД", err)
		return err
	}

	// Проверка соединения с БД
	log.Debugf("Проверка соединения с БД")
	err = db.Ping()
	if err != nil {
		log.Errorf("Ошибка проверки соединения с БД", err)
		return err
	}
	return nil
}

func SetLogFormat() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	//customFormatter.Format()
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(log.DebugLevel)			// Уровень логирования
}
