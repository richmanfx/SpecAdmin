package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"../../models"
)

// Добавляет сценарий в БД
func AddTestScript(newScriptName string, scriptSerialNumber string, scriptSuite string) error {
	var err error
	SetLogFormat()

	// Проверить пермишен пользователя для создания
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "create_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err != nil {
			return err
		}

		// Добавление Скрипта в БД
		log.Debugf("Добавление Скрипта: '%s', Порядковый номер '%s', Сюита: '%s'",
			newScriptName, scriptSerialNumber, scriptSuite)

		result, err := db.Exec("INSERT INTO tests_scripts (name, serial_number, name_suite) VALUES (?,?,?)",
			newScriptName, scriptSerialNumber, scriptSuite)
		if err == nil {
			affected, err := result.RowsAffected()
			if err != nil {
				panic(err)
			}
			log.Debugf("Вставлено %d строк в таблицу 'tests_scripts'.", affected)
		}
		db.Close()
	}
	return err
}


// Удаляет сценарий из БД
func DelTestScript(scriptName, scriptsSuiteName string) error {
	var err error
	SetLogFormat()

	// Проверить пермишен пользователя для удалений
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "delete_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err != nil {
			return err
		}

		// Удаление скрипта из БД
		log.Debugf("Удаление Скрипта '%s'.", scriptName)
		result, err := db.Exec("DELETE FROM tests_scripts WHERE name=? AND name_suite=?", scriptName, scriptsSuiteName)
		if err == nil {
			var affected int64
			affected, err = result.RowsAffected()
			if err == nil {
				if affected == 0 {
					_, goModuleName, lineNumber, _ := runtime.Caller(1)
					log.Debugf("Ошибка удаления Скрипта '%s'. goModuleName=%v, lineNumber=%v",
						scriptName, goModuleName, lineNumber)
				}
				log.Debugf("Удалено строк в БД: %v.", affected)
			}
		}
		db.Close()
	}
	return err
}

// Возвращает Сценарий из БД по имени Сценария и имени его Сюиты
func GetScript(scriptsName, scriptsSuiteName string) (models.Script, int, error) {
	var err error
	var script models.Script
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return script, 0, err }

	// Получить данные о Сценарии и его ключ 'id'
	log.Debugf("Получение данных Сценария '%s' из БД.", scriptsName)
	rows := db.QueryRow("SELECT id,serial_number FROM tests_scripts WHERE name=? AND name_suite=?",
		scriptsName, scriptsSuiteName)

	// Получить данные из результата запроса
	var id int
	var serialNumber string

	err = rows.Scan(&id, &serialNumber)
	if err != nil {
		log.Errorf("Ошибка при получении данных по Сценарию '%s' Сюиты '%s' из БД.", scriptsName, scriptsSuiteName)
	} else {
		log.Debugf("rows.Next из таблицы tests_scripts: %d, %s, %s", id, serialNumber, scriptsSuiteName)

		// Заполнить данными Сценарий
		script.Name = scriptsName
		script.SerialNumber = serialNumber
		script.Suite = scriptsSuiteName
	}

	db.Close()
	log.Debugf("Получение данных Сценария из БД => ошибка '%v'.", err)
	return script, id, err
}

// Обновить данные Сценария в БД
func UpdateTestScript(scriptId int, scriptName string, scriptSerialNumber string, scriptsSuite string) error {
	var err error
	SetLogFormat()

	// Проверить пермишен пользователя для редактирования
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "edit_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err != nil {
			return err
		}

		// Обновить данные о Сценарии в БД
		log.Debugf("Обновление данных о Сценарии '%s' в БД", scriptName)
		_, err = db.Query("UPDATE tests_scripts SET name=?, serial_number=?, name_suite=? WHERE id=? LIMIT 1",
			scriptName, scriptSerialNumber, scriptsSuite, scriptId)
		if err == nil {
			log.Debugf("Успешно обновлены данные Сценария '%s' в БД.", scriptName)
		} else {
			log.Debugf("Ошибка обновления данных Сценария '%s' в БД.", scriptName)
		}
		db.Close()
	}
	return err
}


// Получить список имён Сюит в заданной Группе
func GetSuitsNameFromSpecifiedGroup(groupName string) ([]string, error) {
	var err error
	suitsNameList := make([]string, 0, 0)
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return suitsNameList, err }

	// Запрос имён Сюит
	rows, err := db.Query("SELECT name FROM tests_suits WHERE name_group=? ORDER BY serial_number", groupName)
	if err != nil {panic(err)}		// TODO: Обработать и вывести в браузер
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {	panic(err) }	// TODO: Обработать и вывести в браузер
		log.Debugf("rows.Next из таблицы tests_suits: %s",name)
		suitsNameList = append(suitsNameList, name)		// Добавить в список
	}
	db.Close()
	return suitsNameList, err
}

// Получить только ID Сценариев для заданных Сюит
func GetScriptIdList(suitsNameFromGroup []string) ([]int, error) {

	var err error
	scriptsIdList := make([]int, 0, 0)

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		for _, suiteName := range suitsNameFromGroup {
			rows, _ := db.Query("") // Костыль - пока не понял как объявить отдельно
			rows, err = db.Query("SELECT id FROM tests_scripts WHERE name_suite=? ORDER BY serial_number", suiteName)
			if err == nil {

				// Данные получить из результата запроса
				for rows.Next() {
					var id int
					err = rows.Scan(&id)
					if err == nil {
						log.Debugf("rows.Next из таблицы tests_scripts: %d", id)

						// Добавить в список имён
						scriptsIdList = append(scriptsIdList, id)
					}
				}
			}
		}

	}
	return scriptsIdList, err
}


// Получить Сценарии только для заданных Сюит
func GetScriptListForSpecifiedSuits(suitsNameFromGroup []string) ([]models.Script, error) {
	var err error
	var stepsList []models.Step
	scriptsList := make([]models.Script, 0, 0)

	// Получить только ID Сценариев для заданных Сюит
	scriptsIdList, err := GetScriptIdList(suitsNameFromGroup)
	if err == nil {

		// Получить Шаги из БД только для заданных по ID Сценариев
		stepsList, err = GetStepsListForSpecifiedScripts(scriptsIdList)
		if err == nil {

			// Запрос Сценариев заданных Сюит из БД
			for _, suiteName := range suitsNameFromGroup {
				rows, _ := db.Query("")		// Костыль - пока не понял как обявить отдельно
				rows, err = db.Query(
					"SELECT id, name, serial_number, name_suite FROM tests_scripts WHERE name_suite=? ORDER BY serial_number", suiteName)
				if err == nil {
					// Данные получить из результата запроса
					for rows.Next() {
						var script models.Script
						err = rows.Scan(&script.Id, &script.Name, &script.SerialNumber, &script.Suite)
						if err == nil {

							log.Debugf("rows.Next из таблицы tests_scripts: %s, %s, %s, %s", script.Id, script.Name, script.SerialNumber, script.Suite)

							// Закинуть Шаги в Сценарий
							for _, step := range stepsList {
								if step.ScriptsId == script.Id { // Если Шаг принадлежит Сценарию, то добавляем его
									script.Steps = append(script.Steps, step)
									log.Debugf("Добавлен шаг '%v' в сценарий '%v'", step.Name, script.Name)
								} else {
									log.Debugf("Не добавлен шаг '%v' в сценарий '%v'", step.Name, script.Name)
								}
							}

							scriptsList = append(scriptsList, script) // Добавить сценарий в список
						}
						log.Debugf("Список сценариев: %v", scriptsList)
					}
				}
			}
		}
	}
	return scriptsList, err
}


// Получить список ВСЕХ сценариев
func GetScriptList() ([]models.Script, error)  {

	scriptsList := make([]models.Script, 0, 0)

	// Получить все Шаги из БД
	stepsList := make([]models.Step, 0, 0) // Слайс из Шагов
	stepsList, err := GetStepsList(stepsList)

	// Запрос всех Сценариев из БД
	rows, err := db.Query("SELECT id, name, serial_number, name_suite FROM tests_scripts ORDER BY serial_number")
	if err != nil {panic(err)}		// TODO: Обработать и вывести в браузер

	// Данные получить из результата запроса
	for rows.Next() {
		var id int
		var name string
		var serial_number string
		var name_suite string
		err = rows.Scan(&id, &name, &serial_number, &name_suite)
		if err != nil {
			panic(err)			// TODO: Обработать и вывести в браузер
		}
		log.Debugf("rows.Next из таблицы tests_scripts: %s, %s, %s, %s", id, name, serial_number, name_suite)

		// Заполнить Сценариями список Сценариев
		var script models.Script
		script.Id = id
		script.Name = name
		script.SerialNumber = serial_number
		script.Suite = name_suite
		scriptsList = append(scriptsList, script)
		log.Debugf("Список сценариев: %v", scriptsList)
	}
	return scriptsList, err
}


// Вернуть Сценарий и Сюиту Шага по его ID
func GetScriptAndSuiteByScriptId(ScriptId int) (string, string, error) {
	var err error
	var stepsScriptName string
	var scripsSuiteName string
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return stepsScriptName, scripsSuiteName, err }

	// Данные по Сценарию из БД
	log.Debugf("Получение данных из БД по Сценарию с Id='%s'.", ScriptId)
	rows := db.QueryRow("SELECT name, name_suite FROM tests_scripts WHERE id=?", ScriptId)

	// Получить данные из результата запроса
	err = rows.Scan(&stepsScriptName, &scripsSuiteName)
	if err != nil {
		log.Debugf("Ошибка при получении данных из БД по Сценарию с Id='%s'.", ScriptId)
	} else {
		log.Debugf("rows.Next из таблицы tests_scripts: %s, %s", stepsScriptName, scripsSuiteName)
	}
	return stepsScriptName, scripsSuiteName, err
}
