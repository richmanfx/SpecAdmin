package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	//"fmt"
	//"../../models"
	//"runtime"
)

func AddTestScript(newScript string, scriptSerialNumber string, scriptSuite string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Добавление Скрипта в базу
	log.Debugf("Добавление Скрипта: %s, Порядковый номер '%s' Сюита: %s",
		newScript, scriptSerialNumber, scriptSuite)

	result, err := db.Exec("INSERT INTO tests_scripts (name, serial_number, name_suite) VALUES (?,?,?)",
		newScript, scriptSerialNumber, scriptSuite)
	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {panic(err)}
		log.Debugf("Вставлено строк: %v", affected)
	}
	db.Close()
	return err
}
