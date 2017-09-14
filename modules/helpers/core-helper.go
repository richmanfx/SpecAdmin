package helpers

import (
	"database/sql"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"../../models"
	"fmt"
	"crypto/md5"
	"crypto/sha512"
	"io"
	"time"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

var db *sql.DB

// Соединение с БД и проверка соединения
func dbConnect() error {
	var err error
	SetLogFormat()

	// Проверка соединения с БД
	log.Debugf("Состояние соединения с БД перед подключением => db: '%v'", db)

	// Соединение с БД
	log.Debugf("Подключение к БД")
	db, err = sql.Open("mysql", "specuser:Ghashiz7@tcp(localhost:3306)/specadmin?charset=utf8&parseTime=true")
	if err == nil {
		// Проверка соединения с БД
		log.Debugf("Проверка соединения с БД после ")
		err = db.Ping()
	}
	if err != nil { log.Errorf("Ошибка подключения к БД: '%v'", err) }
	return err
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
	if err == nil {

		// Количество Групп тестов
		requestResult := db.QueryRow("SELECT COUNT(*) FROM tests_groups")
		err = requestResult.Scan(&statistic.GroupsQuantity)
		if err == nil {
			log.Debugf("Количество Групп: '%d'", statistic.GroupsQuantity)
			if statistic.GroupsQuantity == 0 {
				log.Errorln("Нет Групп в таблице 'tests_groups'")
			}

			// Количество Сюит тестов
			requestResult = db.QueryRow("SELECT COUNT(*) FROM tests_suits")
			err = requestResult.Scan(&statistic.SuitesQuantity)
			if err == nil {
				log.Debugf("Количество Сюит: '%d'", statistic.SuitesQuantity)
				if statistic.SuitesQuantity == 0 {
					log.Errorln("Нет Сюит в таблице 'tests_suits'")
				}

				// Количество Сценариев
				requestResult = db.QueryRow("SELECT COUNT(*) FROM tests_scripts")
				err = requestResult.Scan(&statistic.ScriptsQuantity)
				if err == nil {
					log.Debugf("Количество Сценариев: '%d'", statistic.ScriptsQuantity)
					if statistic.ScriptsQuantity == 0 {
						log.Errorln("Нет Сценариев в таблице 'tests_scripts'")
					}

					// Количество Шагов
					requestResult = db.QueryRow("SELECT COUNT(*) FROM tests_steps")
					err = requestResult.Scan(&statistic.StepsQuantity)
					if err == nil {
						log.Debugf("Количество Шагов: '%d'", statistic.StepsQuantity)
						if statistic.StepsQuantity == 0 {
							log.Errorln("Нет Шагов в таблице 'tests_steps'")
						}
					}
				}
			}
		}
	}
	defer db.Close()
	if err != nil { log.Errorf("Ошибка при получении статистики: '%v'", err) }
	return statistic, err
}


// Сгенерировать уникальную строку в 32 hex-символа
func GetUnique32SymbolsString() string {
	h := md5.New()
	io.WriteString(h, time.Now().String())
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Сгенерировать "соль"
func CreateSalt() string {
	hash := sha512.New()
	io.WriteString(hash, time.Now().String())
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Получить Хеш пароля с заданной солью
func CreateHash(password string, salt string) string {
	intSalt, _ := strconv.Atoi(salt)
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), intSalt)
	return string(hashedPassword)
}

// Получить "соль" из БД для заданного пользователя
func GetSaltFromDb(userName string) (string, error) {

	var err error
	var salt string
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {
		// Получить "соль"
		err = db.QueryRow("SELECT salt FROM user WHERE login=?", userName).Scan(&salt)
	}
	defer db.Close()
	if err != nil { log.Errorf("Ошибка получения 'соли' для пользователя с логином '%s': %v", userName, err) }
	return salt, err
}


func ConnectToDB() error {
	return dbConnect()
}


func CloseConnectToDB() error {
	return db.Close()
}


// Получить хеш из БД для заданного пользователя
func GetHashFromDb(userName string) (string, error) {

	var err error
	var hash string
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {
		// Получить "Хеш"
		requestResult := db.QueryRow("SELECT passwd FROM user WHERE login=?", userName)
		err = requestResult.Scan(&hash)
	}
	defer db.Close()
	if err != nil { log.Errorf("Ошибка получения из базы Хеша пароля для пользователя с логином '%s': %v", userName, err) }
	return hash, err
}
