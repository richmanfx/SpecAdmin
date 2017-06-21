package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"runtime"
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
	log.Infof("Удаление Скрипта '%s'.", scriptName)
	result, err := db.Exec("DELETE FROM tests_scripts WHERE name=?", scriptName)
	if err == nil {
		var affected int64
		affected, err = result.RowsAffected()
		if err == nil {
			if affected == 0 {
				_, goFunctionName, lineNumber, _ := runtime.Caller(1)
				err = fmt.Errorf("Ошибка удаления Скрипта '%s'. Есть такой Скрипт?", scriptName)
				log.Infof("Ошибка удаления Скрипта '%s'. goFunctionName=%v, lineNumber=%v",
					scriptName, goFunctionName, lineNumber)
			}
			log.Infof("Удалено строк в БД: %v.", affected)
		}
	}
	
	db.Close()
	return err
}
