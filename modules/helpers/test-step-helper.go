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
	log.Debugf("Добавление Шага '%s', Порядковый номер '%s', Сценарий '%s', Сюита: '%s'.",
		newStepName, stepSerialNumber, stepsScriptName, scriptsSuiteName)

	// Получить ID Сценария в Сюите
	requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?", stepsScriptName, scriptsSuiteName)
	log.Debugf("requestResult: %v", requestResult)

	// Получить значение ID Сценария
	var id int
	err = requestResult.Scan(&id)
	log.Debugf("ID Сценария: '%d'", id)
	if id == 0 {
		log.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.", stepsScriptName, scriptsSuiteName)
		err = fmt.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.", stepsScriptName, scriptsSuiteName)
	} else {
		// В основную таблицу с Шагами
		_, err := db.Exec(
			"INSERT INTO tests_steps (name, serial_number, description, expected_result, script_id) VALUES (?,?,?,?,?)",
								newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult, id)
		if err != nil {panic(err)}		// TODO: Сделать обработку и в Браузер
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

	// Удаление Шага из БД
	log.Debugf("Удаление Шага '%s' из Сценария '%s' Сюиты '%s'.", deletedStepName, stepsScriptName, scriptsSuiteName)

	// Получить ID Сценария в Сюите
	requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?", stepsScriptName, scriptsSuiteName)
	log.Debugf("requestResult: %v", requestResult)

	// Получить значение ID Сценария
	var scriptId int
	err = requestResult.Scan(&scriptId)
	log.Debugf("ID Сценария: '%d'", scriptId)
	if scriptId == 0 {
		log.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.", stepsScriptName, scriptsSuiteName)
		err = fmt.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.", stepsScriptName, scriptsSuiteName)
	} else {
		log.Debugf("Удаление Шага '%s' из Сценария '%s' Сюиты '%s'.", deletedStepName, stepsScriptName, scriptsSuiteName)
		result, err := db.Exec("DELETE FROM tests_steps WHERE name=? AND script_id=?", deletedStepName, scriptId)
		if err != nil {
			err = fmt.Errorf("Ошибка удаления Шага '%s'. Есть такой Шаг?", deletedStepName)
			log.Debugf("Ошибка удаления Шага '%s'", deletedStepName)
			return err
		} else {
			var affected int64
			affected, err = result.RowsAffected()

			if err == nil {
				if affected == 0 {
					err = fmt.Errorf("Ошибка удаления Шага '%s'. Есть такой Шаг?", deletedStepName)
					log.Debugf("Ошибка удаления Шага '%s'", deletedStepName)
					return err
				}
				log.Debugf("Удалено строк в БД: %d", affected)
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
			"SELECT id,name,serial_number,description,expected_result,screen_shot_file_name,script_id FROM tests_steps WHERE script_id=? ORDER BY serial_number",
			scriptsId)
		if err != nil {	panic(err) }	// TODO: Выводить в браузер ошибку

		// Получить данные из результата запроса
		for rows.Next() {
			var stepsId int
			var stepsName string
			var stepsSerialNumber int
			var stepsDescription string
			var stepsExpectedResult string
			var stepsScreenShotsFileName string
			var stepsScriptId int
			err = rows.Scan(&stepsId, &stepsName, &stepsSerialNumber, &stepsDescription, &stepsExpectedResult, &stepsScreenShotsFileName, &stepsScriptId)
			if err != nil {panic(err)}
			log.Debugf("rows.Next из таблицы tests_steps: %d, %s, %d, %s, %s, %s, %d",
				stepsId, stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult, stepsScreenShotsFileName, stepsScriptId)

			// Заполнить Шагами список Шагов
			var step models.Step
			step.Id = stepsId
			step.Name = stepsName
			step.SerialNumber = stepsSerialNumber
			step.Description = stepsDescription
			step.ExpectedResult = stepsExpectedResult
			step.ScreenShotFileName = stepsScreenShotsFileName
			step.ScriptsId = stepsScriptId
			stepsList = append(stepsList, step)
		}
		log.Debugf("Список Шагов: %v", stepsList)
	}

	return stepsList, err
}

// Возвращает Шаг из БД
func GetStep(editedStepName, stepsScriptName, scriptsSuiteName string) (models.Step, error) {
	var err error
	var step models.Step
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return step, err }

	// Получить данные о Шаге
	log.Infof("Получение данных о Шаге '%s' из Сценария '%s' Сюиты '%s'.", editedStepName, stepsScriptName, scriptsSuiteName)

	// Получить ID Сценария в Сюите
	requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?", stepsScriptName, scriptsSuiteName)
	log.Infof("requestResult: %v", requestResult)

	var scriptId int
	err = requestResult.Scan(&scriptId)		// Получить значение ID Сценария
	log.Infof("ID Сценария: '%d'", scriptId)
	if scriptId == 0 {
		log.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.", stepsScriptName, scriptsSuiteName)
		err = fmt.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.", stepsScriptName, scriptsSuiteName)
	} else {
		// Получить Шаг из БД
		log.Infof("Получение Шага '%s' в сценарии с ID '%d'.", editedStepName, scriptId)
		requestResult := db.QueryRow("SELECT id,serial_number,description,expected_result,screen_shot_file_name,script_id FROM tests_steps WHERE name=? AND script_id=?", editedStepName, scriptId)

		// Получить результаты запроса
		var stepId int
		var serialNumber int
		var description string
		var expectedResult string
		var scriptsId int
		var screenShotFileName string

		err = requestResult.Scan(&stepId, &serialNumber, &description, &expectedResult, &screenShotFileName, &scriptsId)
		if err != nil {
			log.Infof("Ошибка при получении данных шага '%s' из БД.", editedStepName)
		} else {
			log.Infof("rows.Next из таблицы tests_steps: %d, %d, %s, %s, %d",
				stepId, serialNumber, description, expectedResult, scriptsId)

			// Заполнить данными Шаг
			step.Id = stepId
			step.Name = editedStepName
			step.SerialNumber = serialNumber
			step.Description = description
			step.ExpectedResult = expectedResult
			step.ScriptsId = scriptsId
			step.ScreenShotFileName = screenShotFileName
		}
	}
	db.Close()
	return step, err
}


// Обновить данные Шага в БД
func UpdateTestStep(
	stepsId int, stepsName string, stepsSerialNumber int, stepsDescription string, stepsExpectedResult string, screenShotsFileName string) error {

	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Обновить данные о Шаге в БД
	log.Infof("Обновление данных о Шаге '%s' в БД", stepsName)
	_, err = db.Query("UPDATE tests_steps SET name=?, serial_number=?, description=?, expected_result=?, screen_shot_file_name=? WHERE id=? LIMIT 1",
		stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult, screenShotsFileName, stepsId)
	if err == nil {
		log.Infof("Успешно обновлены данные Шага '%s' в БД.", stepsName)
	} else {
		log.Infof("Ошибка обновления данных Шага '%s' в БД.", stepsName)
	}
	db.Close()
	return err
}
