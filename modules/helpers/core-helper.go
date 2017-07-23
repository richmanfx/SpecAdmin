package helpers

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"../../models"
	"fmt"
	"crypto/md5"
	"io"
	"time"
)

var db *sql.DB

// Соединение с БД и проверка соединения
func dbConnect() error {

	var err error

	SetLogFormat()

	// Соединение с БД
	log.Debugf("Соединение с БД")
	db, err = sql.Open("mysql", "specuser:Ghashiz7@tcp(localhost:3306)/specadmin?charset=utf8&parseTime=true")
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


// Установить формат логов и уровень логирования
func SetLogFormat() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	//customFormatter.Format()
	log.SetFormatter(customFormatter)
	customFormatter.FullTimestamp = true
	log.SetLevel(log.InfoLevel)			// Уровень логирования
}


// Получить статистику
func GetStatistic() (models.Statistic, error) {

	var err error
	var statistic models.Statistic

	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return statistic, err }


	// Количество Групп тестов
	requestResult := db.QueryRow("SELECT COUNT(*) FROM tests_groups")
	err = requestResult.Scan(&statistic.GroupsQuantity)
	log.Debugf("Количество Групп: '%d'", statistic.GroupsQuantity)
	if statistic.GroupsQuantity == 0 {
		log.Errorln("Нет Групп в таблице 'tests_groups'")
	}

	// Количество Сюит тестов
	requestResult = db.QueryRow("SELECT COUNT(*) FROM tests_suits")
	err = requestResult.Scan(&statistic.SuitesQuantity)
	log.Debugf("Количество Сюит: '%d'", statistic.SuitesQuantity)
	if statistic.SuitesQuantity == 0 {
		log.Errorln("Нет Сюит в таблице 'tests_suits'")
	}

	// Количество Сценариев
	requestResult = db.QueryRow("SELECT COUNT(*) FROM tests_scripts")
	err = requestResult.Scan(&statistic.ScriptsQuantity)
	log.Debugf("Количество Сценариев: '%d'", statistic.ScriptsQuantity)
	if statistic.ScriptsQuantity == 0 {
		log.Errorln("Нет Сценариев в таблице 'tests_scripts'")
	}

	// Количество Шагов
	requestResult = db.QueryRow("SELECT COUNT(*) FROM tests_steps")
	err = requestResult.Scan(&statistic.StepsQuantity)
	log.Debugf("Количество Шагов: '%d'", statistic.StepsQuantity)
	if statistic.StepsQuantity == 0 {
		log.Errorln("Нет Шагов в таблице 'tests_steps'")
	}

	db.Close()
	return statistic, err
}


// Получить уникальное имя файла - 32 hex-символа
func GetUniqueFileName() string {
	h := md5.New()
	io.WriteString(h, time.Now().String())
	return fmt.Sprintf("%x", h.Sum(nil))
}

