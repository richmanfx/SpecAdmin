package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"../../models"
)

func GetGroupsList(groupList []models.Group) ([]models.Group, error) {
	var err error

	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return groupList, err }


	// Считать Шаги из БД


	// Считать Сценарии из БД
	scriptList := make([]models.Script, 0, 0)	// Слайс из Сценариев
	// Запрос всех Сценариев из БД
	rows, err := db.Query("SELECT name, serial_number, name_suite FROM tests_scripts  ORDER BY serial_number")
	if err != nil {panic(err)}
	// Данные получить из результата запроса
	for rows.Next() {
		var name string
		var serial_number string
		var name_suite string
		err = rows.Scan(&name, &serial_number, &name_suite)
		if err != nil {panic(err)}
		log.Debugf("rows.Next из таблицы tests_scripts: %s, %s, %s", name, serial_number, name_suite)

		// Заполнить Сценариями список Сценариев
		var script models.Script
		script.Name = name
		script.SerialNumber = serial_number
		script.Suite = name_suite
		scriptList = append(scriptList, script)
	}
	log.Debugf("Список сценариев: %v", scriptList)

	// Считать Сюиты из БД
	suitesList := make([]models.Suite, 0, 0)	// Слайс из Сюит
	// Запрос Сюит из БД, получить записи
	rows, err = db.Query("SELECT name, description, serial_number ,name_group FROM tests_suits ORDER BY serial_number")
	if err != nil {panic(err)}
	// Данные получить из результата запроса
	for rows.Next() {
		var name string
		var description string
		var serial_number string
		var name_group string
		err = rows.Scan(&name, &description, &serial_number, &name_group)
		if err != nil {panic(err)}
		log.Debugf("rows.Next из таблицы tests_suits: %s, %s, %s, %s", name, description, serial_number, name_group)

		// Заполнить Сюитами список Сюит
		var suite models.Suite
		suite.Name = name
		suite.Description = description
		suite.SerialNumber = serial_number
		suite.Group = name_group

		// Закинуть Сценарии в соответствующие Сюиты
		for _, script := range scriptList {		// Бежать по всем сценариям
			if script.Suite == suite.Name {		// Если Сценарий принадлежит Сюите, то добавляем его
				suite.Scripts = append(suite.Scripts, script)
				log.Debugf("Добавлен сценарий '%v'('%v') в сюиту '%v'", script.Name, script.Suite, suite.Name)
			} else {
				log.Debugf("Не добавлен сценарий '%v'('%v') в сюиту '%v'", script.Name, script.Suite, suite.Name)
			}
		}
		suitesList = append(suitesList, suite)
		log.Debugf("Список Сюит: %v", suitesList)
	}

	// Считать группы из БД
	rows, err = db.Query("SELECT name FROM tests_groups ORDER BY name")
	if err != nil {panic(err)}
	// Данные получить из результата запроса
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {panic(err)}
		log.Debugf("rows.Next из таблицы tests_groups: %s", name)

		// Заполнить Группами список Групп
		var group models.Group
		group.Name = name

		// Закинуть Сюиты в соответствующие Группы
		for _, suite := range suitesList {	// Бежать по всем Сюитам
			if suite.Group == group.Name {	// Если Сюита принадлежит Группе, то добавляем её
				group.Suites = append(group.Suites, suite)
			}
		}
		groupList = append(groupList, group)
	}
	return groupList, err
}

// Добавляет в базу новую группу тестов
// Получает имя новой группы
func AddTestGroup(groupName string) error {
	var err error

	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {
		return err
	}

	// Добавление Группы в базу, используем плейсхолдер
	result, err := db.Exec("INSERT INTO tests_groups (name) VALUE (?)", groupName)
	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {panic(err)}
		log.Debugf("Вставлено строк: %v", affected)
	}
	db.Close()
	return err
}

// Удаляет из базы заданную группу
// Получает имя удаляемой группы
func DelTestGroup(groupName string) error {
	var err error

	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Удаление Группы из базы
	log.Debugf("Удаление Группы: %s", groupName)
	result, err := db.Exec("DELETE FROM tests_groups WHERE name=?", groupName)
	if err == nil {
		var affected int64
		affected, err = result.RowsAffected()
		if err == nil {
			if affected == 0 {
				err = fmt.Errorf("Ошибка удаления Группы '%s'. Есть такая Группа?", groupName)
				log.Debugf("Ошибка удаления Группы '%s'", groupName)
			}
			log.Debugf("Удалено строк в БД: %v", affected)
		}
	}
	db.Close()
	return err
}

// Изменяет имя группы тестов
// Получает имя редактируемой группы
func EditTestGroup(oldName string, newName string) error {
	var err error

	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Изменение имени Группы
	result, err := db.Exec("UPDATE tests_groups SET name=? WHERE name=?", newName, oldName)
	if err != nil {panic(err)}

	affected, err := result.RowsAffected()
	if err != nil {panic(err)}
	log.Debugf("Изменено строк: %v", affected)

	if affected == 0 {
		err = fmt.Errorf("Ошибка изменения имени группы '%s' на новое имя '%s'", oldName, newName)
		log.Debugf("Ошибка изменения имени группы '%s' на новое имя '%s'", oldName, newName )
	}

	db.Close()
	return err
}

