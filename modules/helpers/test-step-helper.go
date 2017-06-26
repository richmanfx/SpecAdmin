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

	// Получить ID из промежуточной таблицы, соответствующий связке Сюиты и Сценария
	requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?", stepsScriptName, scriptsSuiteName)
	log.Infof("requestResult: %v", requestResult)

	// Получить ID сввязки Сценария и Сюиты
	var id int
	err = requestResult.Scan(&id)
	log.Infof("ID: %d", id)
	if id == 0 {
		log.Infof("Не найдено связки Сюиты '%s' и Сценария '%s' в таблице 'tests_scripts'.", scriptsSuiteName, stepsScriptName)
	} else {
		// В основную таблицу с Шагами
		result1, err := db.Exec("INSERT INTO tests_steps (name, serial_number, description, expected_result) VALUES (?,?,?,?)",
								newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult)
		if err != nil {panic(err)}

		// Получить ID новой записи в таблице tests_steps
		stepsId, err := result1.LastInsertId()
		if err != nil {panic(err)}

		// В промежуточную таблицу с ID-шниками
		result2, err := db.Exec("INSERT INTO intermediate_scripts_steps (scripts_id, steps_id) VALUES (?,?)", id, stepsId)
		if err == nil {
			affected, err := result2.RowsAffected()
			if err != nil {panic(err)}
			log.Infof("Вставлено %d строк в таблицу 'intermediate_scripts_steps'.", affected)
		}
	}
	db.Close()
	return err
}
