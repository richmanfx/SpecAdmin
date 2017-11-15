package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"../../models"
	"os"
	"database/sql"
	"errors"
)

// Добавить Шаг в БД
func AddTestStep(newStepName string, stepSerialNumber string, stepsDescription string,stepsExpectedResult string,
				 stepsScreenShotFileName string,	stepsScriptName string, scriptsSuiteName string) error {
	var err error
	SetLogFormat()

	// Проверить пермишен пользователя для создания
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "create_permission")

	if err == nil {
		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Добавление Шага в БД
			log.Debugf("Добавление Шага '%s', Порядковый номер '%s', Сценарий '%s', Сюита: '%s'.",
				newStepName, stepSerialNumber, stepsScriptName, scriptsSuiteName)

			// Получить ID Сценария в Сюите
			requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?",
				stepsScriptName, scriptsSuiteName)
			log.Debugf("requestResult: %v", requestResult)

			// Получить значение ID Сценария
			var id int
			err = requestResult.Scan(&id)
			log.Debugf("ID Сценария: '%d'", id)
			if id == 0 {
				log.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.",
							stepsScriptName, scriptsSuiteName)
				err = errors.New(
					fmt.Sprintf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.",
								 stepsScriptName, scriptsSuiteName))
			} else {
				// В таблицу с Шагами
				_, err = db.Exec(
					"INSERT INTO tests_steps (name, serial_number, description, expected_result, screen_shot_file_name, script_id) VALUES (?,?,?,?,?,?)",
					newStepName, stepSerialNumber, stepsDescription, stepsExpectedResult, stepsScreenShotFileName, id)
			}
		}
		defer db.Close()
	}
	if err != nil {log.Errorf("Ошибка при добавлении Шага '%s' в табицу 'tests_steps': '%v'", newStepName, err)}
	return err
}


// Удалить Шаг из БД
func DelTestStep(deletedStepName, stepsScriptName, scriptsSuiteName string) error {
	var err error
	var screenShotsFileName string // Имя файла скриншота - удалим после удаления Шага из БД
	SetLogFormat()

	// Проверить пермишен пользователя для удалений
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "delete_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Удаление Шага из БД
			log.Debugf("Удаление Шага '%s' из Сценария '%s' Сюиты '%s'.",
					    deletedStepName, stepsScriptName, scriptsSuiteName)

			// Получить ID Сценария в Сюите
			requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?",
										  stepsScriptName, scriptsSuiteName)
			log.Debugf("requestResult: %v", requestResult)

			// Получить значение ID Сценария
			var scriptId int
			err = requestResult.Scan(&scriptId)
			log.Debugf("ID Сценария: '%d'", scriptId)
			if scriptId > 0 {

				// Получить имя файла скриншота Шага
				var result sql.Result
				screenShotsFileName, err = GetStepsScreenShotsFileName(deletedStepName, scriptId)
				if err == nil {

					log.Debugf("Удаление Шага '%s' из Сценария '%s' Сюиты '%s'.",
								deletedStepName, stepsScriptName, scriptsSuiteName)
					result, err = db.Exec("DELETE FROM tests_steps WHERE name=? AND script_id=?",
										   deletedStepName, scriptId)
					if err == nil {

						var affected int64
						affected, err = result.RowsAffected()

						if err == nil {
							if affected != 0 {
								log.Debugf("Удалено строк в БД: %d", affected)
							} else {
								log.Errorf("Ошибка удаления Шага '%s'", deletedStepName)
							}
						}
					}
				}
			} else {
				log.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.",
							stepsScriptName, scriptsSuiteName)
				err = errors.New(
					fmt.Sprintf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.",
								 stepsScriptName, scriptsSuiteName))
			}
		}
		defer db.Close()

		// Удаление файла скриншота из хранилища
		if err == nil && screenShotsFileName != "" { // Если из базы удалили без ошибок
			err = DelScreenShotsFile(screenShotsFileName)

			if err != nil {	log.Errorf("Ошибка удаления файла скриншота '%s'", screenShotsFileName)	}
		}
	}
	if err != nil {log.Errorf("Ошибка при удалении Шага из БД: '%v'", err)}
	return err
}


// Вернуть имя файла скриншота Шага по заданным Имени Шага и Id его сценария
func GetStepsScreenShotsFileName(deletedStepName string, scriptId int) (string, error) {
	var screenShotsFileName string

	requestResult := db.QueryRow("SELECT screen_shot_file_name FROM tests_steps WHERE name=? AND script_id=?",
								  deletedStepName, scriptId)
	log.Debugf("requestResult: %v", requestResult)
	err := requestResult.Scan(&screenShotsFileName)		// Получить имя файла скриншота
	log.Debugf("Имя файла скриншота: '%s'", screenShotsFileName)
	if err != nil {
		log.Errorf("Не найдено имя файла скриншота для Шага '%s' в Сценарии с Id '%s' в таблице 'tests_steps'.",
				    deletedStepName, scriptId)
		}
	return screenShotsFileName, err
}


// Сформировать список ВСЕХ Шагов из БД
func GetStepsList(stepsList []models.Step) ([]models.Step, error)  {

	// Запрос всех Шагов из БД
	rows, err := db.Query("SELECT id,name,serial_number,description,expected_result FROM tests_steps ORDER BY serial_number")
	if err == nil {
		// Получить данные из результата запроса
		var step models.Step
		for rows.Next() {
			err = rows.Scan(&step.Id, &step.Name, &step.SerialNumber, &step.Description, &step.ExpectedResult)
			if err == nil {
				log.Debugf("rows.Next из таблицы tests_steps: %s, %s, %d, %s, %s",
					step.Id, step.Name, step.SerialNumber, step.Description, step.ExpectedResult)
				stepsList = append(stepsList, step)
			}
		}
		log.Debugf("Список Шагов: %v", stepsList)
	}
	if err != nil {log.Errorf("Ошибка при формировании всех шагов из БД: '%v'", err)}
	return stepsList, err
}


// Получить Шаги из БД только для заданных по ID Сценариев
func GetStepsListForSpecifiedScripts(scriptsIdList []int) ([]models.Step, error) {
	var err error
	stepsList := make([]models.Step, 0, 0) // Слайс из Шагов

	// Подключиться к БД
	err = dbConnect()
	if err == nil {
		for _, scriptsId := range scriptsIdList {
			rows, err := db.Query(
				"SELECT id,name,serial_number,description,expected_result,screen_shot_file_name,script_id FROM tests_steps WHERE script_id=? ORDER BY serial_number",
				scriptsId)
			if err == nil {

				// Получить данные из результата запроса
				var step models.Step
				for rows.Next() {
					err = rows.Scan(&step.Id, &step.Name, &step.SerialNumber, &step.Description, &step.ExpectedResult,
									&step.ScreenShotFileName, &step.ScriptsId)
					if err == nil {
						log.Debugf("rows.Next из таблицы tests_steps: %d, %s, %d, %s, %s, %s, %d",
								    step.Id, step.Name, step.SerialNumber, step.Description,
								    step.ExpectedResult, step.ScreenShotFileName, step.ScriptsId)
						stepsList = append(stepsList, step) // Добавить шаг в список
					}
				}
				log.Debugf("Список Шагов: %v", stepsList)
			}
		}
	}
	defer db.Close()
	if err != nil {log.Errorf("Ошибка при получении Шагов из БД только для заданных по ID Сценариев: '%v'", err)}
	return stepsList, err
}

// Возвращает Шаг из БД
func GetStep(editedStepName, stepsScriptName, scriptsSuiteName string) (models.Step, error) {
	var err error
	var step models.Step
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {
		// Получить данные о Шаге
		log.Debugf("Получение данных о Шаге '%s' из Сценария '%s' Сюиты '%s'.",
					editedStepName, stepsScriptName, scriptsSuiteName)

		// Получить ID Сценария в Сюите
		requestResult := db.QueryRow("SELECT id FROM tests_scripts WHERE name=? AND name_suite=?",
									  stepsScriptName, scriptsSuiteName)
		log.Debugf("requestResult: %v", requestResult)

		var scriptId int
		err = requestResult.Scan(&scriptId) // Получить значение ID Сценария
		log.Debugf("ID Сценария: '%d'", scriptId)
		if scriptId == 0 {
			log.Errorf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.",
						stepsScriptName, scriptsSuiteName)
			err = errors.New(fmt.Sprintf("Не найдено Сценария '%s' в Сюите '%s' в таблице 'tests_scripts'.",
										  stepsScriptName, scriptsSuiteName))
		} else {
			step.Name = editedStepName

			// Получить Шаг из БД
			log.Debugf("Получение Шага '%s' в сценарии с ID '%d'.", editedStepName, scriptId)
			requestResult := db.QueryRow("SELECT id,serial_number,description,expected_result,screen_shot_file_name,script_id FROM tests_steps WHERE name=? AND script_id=?",
										  editedStepName, scriptId)

			// Получить результаты запроса
			err = requestResult.Scan(&step.Id, &step.SerialNumber, &step.Description, &step.ExpectedResult,
									 &step.ScreenShotFileName, &step.ScriptsId)
			if err == nil {
				log.Debugf("rows.Next из таблицы tests_steps: %d, %d, %s, %s, %d",
					step.Id, step.SerialNumber, step.Description, step.ExpectedResult, step.ScriptsId)
			}
		}
	}
	defer db.Close()
	if err != nil {log.Errorf("Ошибка при получении данных шага '%s' из БД: '%v'", editedStepName, err)}
	return step, err
}


// Обновить данные Шага в БД
func UpdateTestStep(stepsId int, stepsName string, stepsSerialNumber int, stepsDescription string,
					stepsExpectedResult string, screenShotsFileName string) error {
	var err error
	SetLogFormat()

	// Проверить пермишен пользователя для редактирования
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "edit_permission")

	if err == nil {
		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Обновить данные о Шаге в БД
			log.Debugf("Обновление данных о Шаге '%s' в БД", stepsName)

			if screenShotsFileName == "" {
				_, err = db.Exec("UPDATE tests_steps SET name=?, serial_number=?, description=?, expected_result=? WHERE id=? LIMIT 1",
					stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult, stepsId)
			} else {
				_, err = db.Exec("UPDATE tests_steps SET name=?, serial_number=?, description=?, expected_result=?, screen_shot_file_name=? WHERE id=? LIMIT 1",
					stepsName, stepsSerialNumber, stepsDescription, stepsExpectedResult, screenShotsFileName, stepsId)
			}
			if err == nil {
				log.Debugf("Успешно обновлены данные Шага '%s' в БД.", stepsName)
			}
		}
		defer db.Close()
	}
	if err != nil {log.Errorf("Ошибка обновления данных Шага '%s' в БД.", stepsName)}
	return err
}


// Удалить скриншот у Шага по заданному Id Шага
func DeleteStepsScreenShot(stepsId int) error {

	var err error
	var requestResult *sql.Row
	SetLogFormat()

	// Проверить пермишен пользователя для удалений
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "delete_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {
			/// Удалить файл скриншота на диске
			// Получить имя файла скриншота
			requestResult = db.QueryRow("SELECT screen_shot_file_name FROM tests_steps WHERE id=?", stepsId)
		}
		defer db.Close()

		// Получить результаты запроса
		var screenShotFileName string
		err = requestResult.Scan(&screenShotFileName)

		if err != nil {
			log.Errorf("Ошибка при получении имени файла скриншота '%s' из БД.", screenShotFileName)
		} else {
			log.Debugf("rows.Next из таблицы tests_steps: %s", screenShotFileName)
		}

		log.Debugf("Удаление файла скриншота '%s'", screenShotFileName)
		err = DelScreenShotsFile(screenShotFileName)

		if err == nil {
			// Подключиться к БД
			err = dbConnect()
			if err == nil {
				log.Debugln("Подключились к БД.")

				// Удалить данные о скриншоте в БД
				log.Debugf("Удаление данные о скриншоте в БД, Id Шага: '%d'.", stepsId)
				_, err = db.Exec("UPDATE tests_steps SET screen_shot_file_name=? WHERE id=? LIMIT 1", "", stepsId)
			}
			defer db.Close()
		}
		if err == nil {	log.Debugf("Скриншот успешно удалён из БД в Шаге с Id='%d'", stepsId) }
	}
	if err != nil {log.Errorf("Ошибка при удалении скриншота из БД в Шаге с Id='%d', %v", stepsId, err)}
	return err
}


// Удаление файла скриншота по имени файла
func DelScreenShotsFile(screenShotsFileName string) error {

	var fullScreenShotsPath string
	var screenShotsPath string
	var err error

	// Получить путь до хранилища скриншотов
	screenShotsPath, err = GetScreenShotsPath()
	if err == nil {
		lastSymbolOfPath := screenShotsPath[len(screenShotsPath)-1:]
		log.Debugf("Последний символ в пути: '%s'", lastSymbolOfPath)

		if lastSymbolOfPath != string(os.PathSeparator) {
			fullScreenShotsPath = screenShotsPath + string(os.PathSeparator) + screenShotsFileName
		} else {
			fullScreenShotsPath = screenShotsPath + screenShotsFileName
		}

		// Удалить файл
		err = os.Remove(fullScreenShotsPath)
	}
	if err != nil {log.Errorf("Ошибка при удалении файла скриншота по имени файла: %v", err)}
	return err
}

// Возвращает имя, описание и ожидаемый результат шага по его ID
func GetStepsData(stepsId int) (string, string, string, error) {

	var err error
	var requestResult *sql.Row
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {
		// Получить данные
		requestResult = db.QueryRow("SELECT name, description, expected_result FROM tests_steps WHERE id=?", stepsId)
	}
	defer db.Close()

	// Получить результаты запроса
	var stepsName string
	var stepsDescription string
	var stepsExpectedResult string

	err = requestResult.Scan(&stepsName, &stepsDescription, &stepsExpectedResult)

	if err != nil {
		log.Errorf("Ошибка получения данных Шага по ID '%s' из БД.", stepsId)
		}
	return stepsName, stepsDescription, stepsExpectedResult, err
}
