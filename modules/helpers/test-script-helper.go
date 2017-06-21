package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"runtime"
	"../../models"
)

func AddTestScript(newScriptName string, scriptSerialNumber string, scriptSuite string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Добавление Скрипта в БД
	log.Debugf("Добавление Скрипта: %s, Порядковый номер '%s' Сюита: %s",
		newScriptName, scriptSerialNumber, scriptSuite)

	result, err := db.Exec("INSERT INTO tests_scripts (name, serial_number, name_suite) VALUES (?,?,?)",
		newScriptName, scriptSerialNumber, scriptSuite)
	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {panic(err)}
		log.Debugf("Вставлено строк: %v.", affected)
	}
	db.Close()
	return err
}

func DelTestScript(scriptName string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Удаление скрипта из БД
	log.Debugf("Удаление Скрипта '%s'.", scriptName)
	result, err := db.Exec("DELETE FROM tests_scripts WHERE name=?", scriptName)
	if err == nil {
		var affected int64
		affected, err = result.RowsAffected()
		if err == nil {
			if affected == 0 {
				_, goModuleName, lineNumber, _ := runtime.Caller(1)
				err = fmt.Errorf("Ошибка удаления Скрипта '%s'. Есть такой Скрипт?", scriptName)
				log.Debugf("Ошибка удаления Скрипта '%s'. goModuleName=%v, lineNumber=%v",
					scriptName, goModuleName, lineNumber)
			}
			log.Debugf("Удалено строк в БД: %v.", affected)
		}
	}
	
	db.Close()
	return err
}

// Возвращает Сценарий из БД
func GetScript(scriptsName string) (models.Script, int, error) {
	var err error
	var script models.Script
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return script, 0, err }

	// Получить данные о Сценарии и его ключ 'id'
	log.Infof("Получение данных Сценария '%s' из БД.", scriptsName)
	rows := db.QueryRow("SELECT id,serial_number,name_suite FROM tests_scripts WHERE name=?", scriptsName)

	var id int
	var serialNumber string
	var suiteName string
	if err == nil {
		// Получить данные из результата запроса
		err = rows.Scan(&id, &serialNumber, &suiteName)
		if err != nil {
			log.Infof("Ошибка при получении данных по сценарию '%s' из БД.", scriptsName)
		} else {
			log.Infof("rows.Next из таблицы tests_scripts: %d, %s, %s", id, serialNumber, suiteName)

			// Заполнить данными Сценарий
			script.Name = scriptsName
			script.SerialNumber = serialNumber
			script.Suite = suiteName
		}
	}
	db.Close()
	log.Infof("Получение данных Сценария из БД => ошибка '%v'.", err)
	return script, id, err
}

// Обновить данные Сценария в БД
func UpdateTestScript(scriptId int, scriptName string, scriptSerialNumber string, scriptsSuite string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Обновить данные о Сценарии в БД
	log.Infof("Обновление данных о Сценарии '%s' в БД", scriptName)
	_, err = db.Query("UPDATE tests_scripts SET name=?, serial_number=?, name_suite=? WHERE id=? LIMIT 1",
		scriptName, scriptSerialNumber, scriptsSuite, scriptId)
	if err == nil {
		log.Infof("Успешно обновлены данные Сценария '%s' в БД.", scriptName)
	} else {
		log.Infof("Ошибка обновления данных Сценария '%s' в БД.", scriptName)
	}

	return err
}
