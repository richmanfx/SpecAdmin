package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

func AddTestStep(
	newStepName string, stepSerialNumber string, stepsDescription string,
	stepsExpectedResult string, stepsScriptName string, scriptsSuiteName string) error {

	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Добавление Шага в БД
	log.Infof("Добавление Шага '%s', Порядковый номер '%s', Сценарий '%s', Сюита: '%s'.",
		newStepName, stepSerialNumber, stepsScriptName, scriptsSuiteName)

	result, err := db.Exec("INSERT INTO tests_steps (name, serial_number, description, expected_result, name_script, name_suite) VALUES (?,?,?,?,?,?) 	",
		newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult, stepsScriptName, scriptsSuiteName)

	// Подзапросы:
	// SELECT * from tests_scripts WHERE name='Получение тестовых данных и подготовка окружения' and name_suite='FnsOtchetnostUL_Dual'
	// SELECT * from tests_scripts WHERE name='Получение тестовых данных и подготовка окружения' and name_suite='FnsOtchetnostUL_Normal'

	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {panic(err)}
		log.Debugf("Вставлено %d строк в таблицу 'tests_steps'.", affected)
	}
	db.Close()
	return err
}
