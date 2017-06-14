package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

//var db *sql.DB

func GetTestGroupsList() ([]string, error) {

	testGroupList := make([]string, 0, 30)
	var err error

	// Подключиться к БД
	err = dbConnect()
	if err != nil {
		return testGroupList, err
	}

	// Запрос Групп из базы, получить записи
	rows, err := db.Query("SELECT (name) FROM tests_groups")
	if err != nil {panic(err)}

	// Данные получить из результата запроса
	for rows.Next() {
		var group string
		err = rows.Scan(&group)
		if err != nil {panic(err)}
		log.Debugf("row.Next group: %s", group)
		testGroupList = append(testGroupList, group)
	}

	rows.Close()
	db.Close()

	return testGroupList, err
}

// Добавляет в базу новую группу тестов
// Получает имя новой группы
func AddTestGroup(groupName string) error {

	var err error

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
		log.Infof("Вставлено строк: %v", affected)
	}

	db.Close()

	return err
}

// Удаляет из базы заданную группу
// Получает имя удаляемой группы
func DelTestGroup(groupName string) error {

	var err error

	// Подключиться к БД
	err = dbConnect()
	if err != nil {
		return err
	}

	// Удаление Группы из базы
	result, err := db.Exec("DELETE FROM tests_groups WHERE name=?", groupName)
	if err != nil {panic(err)}

	affected, err := result.RowsAffected()
	if err != nil {panic(err)}
	log.Infof("Удалено строк: %v", affected)

	db.Close()

	return err
}

// Изменяет имя группы тестов
// Получает имя редактируемой группы
func EditTestGroup(oldName string, newName string) error {

	var err error

	// Подключиться к БД
	err = dbConnect()
	if err != nil {
		return err
	}

	// Изменение имени Группы
	result, err := db.Exec("UPDATE tests_groups SET name=? WHERE name=?", newName, oldName)
	if err != nil {panic(err)}

	affected, err := result.RowsAffected()
	if err != nil {panic(err)}
	log.Infof("Изменено строк: %v", affected)

	if affected == 0 {
		err = fmt.Errorf("Ошибка изменения имени группы '%s' на новое имя '%s'", oldName, newName)
		log.Infof("Ошибка изменения имени группы '%s' на новое имя '%s'", oldName, newName )
	}

	db.Close()

	return err
}

