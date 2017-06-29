package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"../../models"
	"runtime"
)

// Сформировать список Сюит для указанной Группы
func GetSuitesListInGroup(groupName string) ([]models.Suite, error) {
	var err error
	suitesList := make([]models.Suite, 0, 0)	// Слайс из Сюит
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return suitesList, err }

	//// Считать Сценарии из БД
	scriptsList := make([]models.Script, 0, 0)
	scriptsList, err = GetScriptList(scriptsList)

	// Получить только список имён Сюит в данной Группе
	suitsNameFromGroup, err := GetSuitsNameFromSpecifiedGroup(groupName)
	log.Debugf("Сюиты из группы '%s': %v", groupName, suitsNameFromGroup)


	// Считать Сенарии только для заданных Сюит
	//scriptsList := make([]models.Script, 0, 0)
	//scriptsList, err = GetScriptListForSpecifiedSuits(scriptsList, suitsNameFromGroup)

	for _, suiteName := range suitsNameFromGroup {

		// Сюиты из БД по списку имён Сюит
		rows, err := db.Query("SELECT name,description,serial_number FROM tests_suits WHERE name=? ORDER BY serial_number", suiteName)

		// Получить данные из результата запроса
		for rows.Next() {
			var name string
			var description string
			var serial_number string
			err = rows.Scan(&name, &description, &serial_number)
			if err != nil {
				panic(err)
			}
			log.Debugf("Данные из таблицы 'tests_suits': %s, %s, %s, %s", name, description, serial_number)

			// Заполнить Сюитами список Сюит
			var suite models.Suite
			suite.Name = name
			suite.Description = description
			suite.SerialNumber = serial_number
			suite.Group = groupName

			// Закинуть Сценарии в соответствующие Сюиты
			for _, script := range scriptsList { // Бежать по всем сценариям
				if script.Suite == suite.Name { // Если Сценарий принадлежит Сюите, то добавляем его
					suite.Scripts = append(suite.Scripts, script)
					log.Debugf("Добавлен сценарий '%v'('%v') в сюиту '%v'", script.Name, script.Suite, suite.Name)
				} else {
					log.Debugf("Не добавлен сценарий '%v'('%v') в сюиту '%v'", script.Name, script.Suite, suite.Name)
				}
			}

			// Добавить Сюиту в список
			suitesList = append(suitesList, suite)
		}
	}
	log.Debugf("Список Сюит: '%v'", suitesList)
	return suitesList, err
}


// Добавляет в базу новую сюиту тестов
// Получает имя новой сюиты, описание этой сюиты и группу тестов, в которую вносится сюита
func AddTestSuite(suitesName string, suitesDescription string, suitesSerialNumber string, suitesGroup string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Добавление Сюиты в базу, используем плейсхолдер
	log.Debugf("Добавление Сюиты: %s, Описание: %s, Порядковый номер '%s' Группа: %s",
		suitesName, suitesDescription, suitesSerialNumber, suitesGroup)
	result, err := db.Exec("INSERT INTO tests_suits (name, description, serial_number, name_group) VALUES (?,?,?,?)",
		suitesName, suitesDescription, suitesSerialNumber, suitesGroup)
	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {panic(err)}
		log.Debugf("Вставлено %d строк в таблицу 'tests_suits'.", affected)
	}
	db.Close()
	return err
}

// Удаляет из базы сюиту
// Получает имя удаляемой сюиты
func DelTestSuite(suitesName string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Удаление Сюиты из базы данных
	log.Debugf("Удаление Сюиты '%s'.", suitesName)
	result, err := db.Exec("DELETE FROM tests_suits WHERE name=?", suitesName)
	if err == nil {
		var affected int64
		affected, err = result.RowsAffected()
		if err == nil {
			if affected == 0 {
				_, goModuleName, lineNumber, _ := runtime.Caller(1)
				err = fmt.Errorf("Ошибка удаления Сюиты '%s'. Есть такая Сюита?", suitesName)
				log.Debugf("Ошибка удаления Сюиты '%s'. goModuleName=%v, lineNumber=%v",
					suitesName, goModuleName, lineNumber)
			}
			log.Debugf("Удалено строк в БД: %v", affected)
		}
	}
	db.Close()
	return err
}

// Возвращает Сюиту из БД
func GetSuite(suitesName string) (models.Suite, int, error)  {
	var err error
	var suite models.Suite
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return suite, 0, err }

	// Получить данные о Сюите и её ключ 'id'
	log.Debugf("Получение данных Сюиты '%s' из БД.", suitesName)
	rows := db.QueryRow("SELECT id,description,serial_number,name_group FROM tests_suits WHERE name=?", suitesName)

	var id int
	var description string
	var serialNumber string
	var groupName string
	if err == nil {
		// Данные получить из результата запроса
		err = rows.Scan(&id, &description, &serialNumber, &groupName)
		if err != nil {
			log.Debugf("Ошибка при получении данных по сюите '%s' из БД.", suitesName)
		} else {
			log.Debugf("rows.Next из таблицы tests_suits: %d, %s, %s, %s", id, description, serialNumber, groupName)

			// Заполнить данными Сюиту
			suite.Name = suitesName
			suite.Description = description
			suite.SerialNumber = serialNumber
			suite.Group = groupName
		}
	}
	db.Close()
	log.Debugf("Получение данных Сюиты из БД => ошибка '%v'.", err)
	return suite, id, err
}

// Обновляет данные Сюиты в БД
func UpdateTestSuite(suitesId int, suitesName string, suitesDescription string,
					 suitesSerialNumber string, suitesGroup string) error {

	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Обновить данные о Сюите в БД
	log.Debugf("Обновление данных о Сюите '%s' в БД", suitesName)
	_, err = db.Query("UPDATE tests_suits SET name=?, description=?, serial_number=?, name_group=? WHERE id=? LIMIT 1",
		suitesName, suitesDescription, suitesSerialNumber, suitesGroup, suitesId)
	if err == nil {
		log.Debugf("Успешно обновлены данные Сюиты '%s' в БД.", suitesName)
	} else {
		log.Debugf("Ошибка обновления данных Сюиты '%s' в БД.", suitesName)
	}

	return err
}
