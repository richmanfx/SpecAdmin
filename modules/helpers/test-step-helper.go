package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"../../models"
)

// Добавить Шаг в БД
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

	// Получить ID для промежуточной таблицы, соответствующий связке Сюиты и Сценария
	requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?", stepsScriptName, scriptsSuiteName)
	log.Infof("requestResult: %v", requestResult)

	// Получить ID связки Сценария и Сюиты
	var id int
	err = requestResult.Scan(&id)
	log.Infof("ID (Сюита + Сценарий): %d", id)
	if id == 0 {
		log.Errorf("Не найдено связки Сюиты '%s' и Сценария '%s' в таблице 'tests_scripts'.", scriptsSuiteName, stepsScriptName)
		err = fmt.Errorf("Не найдено связки Сюиты '%s' и Сценария '%s' в таблице 'tests_scripts'.", scriptsSuiteName, stepsScriptName)
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


// Удалить Шаг из БД
func DelTestStep(deletedStepName, stepsScriptName, scriptsSuiteName string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Получить ID из промежуточной таблицы, соответствующий связке Сюиты и Сценария
	requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?", stepsScriptName, scriptsSuiteName)

	// Получить ID связки Сценария и Сюиты
	var id int
	err = requestResult.Scan(&id)
	log.Infof("ID (Сюита + Сценарий): %d", id)
	if id == 0 {
		log.Infof("Не найдено связки Сюиты '%s' и Сценария '%s' в таблице 'tests_scripts'.", scriptsSuiteName, stepsScriptName)
		err = fmt.Errorf("Не найдено связки Сюиты '%s' и Сценария '%s' в таблице 'tests_scripts'.", scriptsSuiteName, stepsScriptName)
	} else {
		// Получить ID шага из промежуточной таблицы
		requestResult = db.QueryRow("SELECT steps_id FROM intermediate_scripts_steps WHERE intermediate_scripts_steps.scripts_id=?", id)
		var stepsId int
		err = requestResult.Scan(&stepsId)
		log.Infof("ID шага в промежуточной таблице: %d", stepsId)

		log.Infof("Удаление Шага '%s' из сценария '%s' сюиты '%s'.", deletedStepName, stepsScriptName, scriptsSuiteName)
		// Удаление Шага из таблицы 'tests_steps'
		result, err := db.Exec("DELETE FROM tests_steps WHERE name=? AND id=?", deletedStepName, stepsId)
		if err != nil {

			var affected int64
			affected, err = result.RowsAffected()

			if err == nil {
				if affected == 0 {
					err = fmt.Errorf("Ошибка удаления Шага '%s'. Есть такой Шаг?", deletedStepName)
					log.Infof("Ошибка удаления Шага '%s'", deletedStepName)
				}
				log.Infof("Удалено строк в БД: %d", affected)
			} else {
				return err
			}

		} else {

			var affected int64
			affected, err = result.RowsAffected()

			if err == nil {
				if affected == 0 {
					err = fmt.Errorf("Ошибка удаления Шага '%s'. Есть такой Шаг?", deletedStepName)
					log.Infof("Ошибка удаления Шага '%s'", deletedStepName)
					return err
				}
				log.Infof("Удалено строк в БД: %d", affected)

				// Удаление ID сценария и ID шага из промежуточной таблицы
				result, err := db.Exec("DELETE FROM intermediate_scripts_steps WHERE steps_id=?", stepsId)
				if err != nil {

					var affected int64
					affected, err = result.RowsAffected()

					if err == nil {
						if affected == 0 {
							err = fmt.Errorf("Ошибка удаления ID '%s' Шага '%s' из промежуточной таблицы.", stepsId, deletedStepName)
							log.Infof("Ошибка удаления ID '%s' Шага '%s' из промежуточной таблицы.", stepsId, deletedStepName)
						}
						log.Infof("Удалено строк в БД: %d", affected)
					} else {
						return err
					}

				}

			} else {
				return err
			}

		}
	}
	db.Close()
	return err
}


// Сформировать список ВСЕХ Шагов из БД
func GetStepsList(stepsList []models.Step) ([]models.Step, error)  {

	// Запрос всех Шагов из БД
	rows, err := db.Query("SELECT id,name,serial_number,description,expected_result FROM tests_steps ORDER BY serial_number")
	if err != nil {panic(err)}

	// Получить данные из результата запроса
	for rows.Next() {
		var stepsId int
		var stepsName string
		var stepsSerialNumber int
		var stepsDescription string
		var stepsExpectedResult string
		err = rows.Scan(&stepsId, &stepsName, &stepsSerialNumber, &stepsDescription, &stepsExpectedResult)
		if err != nil {panic(err)}
		log.Debugf("rows.Next из таблицы tests_steps: %s, %s, %d, %s, %s",
			stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult)

		// Заполнить Шагами список Шагов
		var step models.Step
		step.Id = stepsId
		step.Name = stepsName
		step.SerialNumber = stepsSerialNumber
		step.Description = stepsDescription
		step.ExpectedResult = stepsExpectedResult
		stepsList = append(stepsList, step)
	}
	log.Debugf("Список Шагов: %v", stepsList)
	return stepsList, err
}


// Получить Шаги из БД только для заданных по ID Сценариев
func GetStepsListForSpecifiedScripts(scriptsIdList []int) ([]models.Step, error) {
	var err error
	stepsList := make([]models.Step, 0, 0) // Слайс из Шагов

	for _, scriptsId := range scriptsIdList {
		rows, err := db.Query(
			"SELECT id,name,serial_number,description,expected_result,script_id FROM tests_steps WHERE script_id=? ORDER BY serial_number",
			scriptsId)
		if err != nil {	panic(err) }	// TODO: Выводить в браузер ошибку

		// Получить данные из результата запроса
		for rows.Next() {
			var stepsId int
			var stepsName string
			var stepsSerialNumber int
			var stepsDescription string
			var stepsExpectedResult string
			var stepsScriptId int
			err = rows.Scan(&stepsId, &stepsName, &stepsSerialNumber, &stepsDescription, &stepsExpectedResult, &stepsScriptId)
			if err != nil {panic(err)}
			log.Debugf("rows.Next из таблицы tests_steps: %d, %s, %d, %s, %s, %d",
				stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult, stepsScriptId)

			// Заполнить Шагами список Шагов
			var step models.Step
			step.Id = stepsId
			step.Name = stepsName
			step.SerialNumber = stepsSerialNumber
			step.Description = stepsDescription
			step.ExpectedResult = stepsExpectedResult
			step.ScriptsId = stepsScriptId
			stepsList = append(stepsList, step)
		}
		log.Debugf("Список Шагов: %v", stepsList)
	}

	return stepsList, err
}
