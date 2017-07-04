package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"../../models"
)

// Возвращает список Групп из БД
func GetGroupsList(groupList []models.Group) ([]models.Group, error) {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return groupList, err }

	// Считать группы из БД
	rows, err := db.Query("SELECT name FROM tests_groups ORDER BY name")
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
		groupList = append(groupList, group)
	}
	db.Close()
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
		log.Debugf("Вставлено %d строк в таблиц 'tests_groups'.", affected)
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

