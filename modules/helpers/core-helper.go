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

	// Соединение с БД
	log.Infof("Соединение с БД")
	db, err = sql.Open("mysql", "specuser:Ghashiz7@tcp(localhost:3306)/specadmin?charset=utf8")
	if err != nil {
		log.Errorf("Ошибка соединение с БД", err)
		return err
	}

	// Проверка соединения с БД
	log.Infof("Проверка соединения с БД")
	err = db.Ping()
	if err != nil {
		log.Errorf("Ошибка проверки соединения с БД", err)
		return err
	}
	return nil
}
