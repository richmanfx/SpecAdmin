package helpers

import (
	"SpecAdmin/models"
	"database/sql"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

// Возвращает список Групп из БД
func GetGroupsList(groupList *[]models.Group) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Считать группы из БД
		rows, err := db.Query("SELECT name FROM tests_groups ORDER BY name")
		if err == nil {

			// Данные получить из результата запроса
			var group models.Group
			for rows.Next() {
				err = rows.Scan(&group.Name)
				if err == nil {
					log.Debugf("rows.Next из таблицы tests_groups: %s", group.Name)
					*groupList = append(*groupList, group)
				}
			}
		}
	}
	defer CloseConnectToDB()
	if err != nil {
		log.Errorf("Ошибка при получении списка Групп из БД: %v", err)
	}
	return err
}

// Добавляет в базу новую группу тестов
// Получает имя новой группы
func AddTestGroup(groupName string) error {

	var err error
	var result sql.Result
	var affected int64
	SetLogFormat()

	// Проверить пермишен пользователя для создания
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "create_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Добавление Группы в базу, используем плейсхолдер
			result, err = db.Exec("INSERT INTO tests_groups (name) VALUES (?)", groupName)
			if err == nil {
				affected, err = result.RowsAffected()
				if err == nil {
					log.Debugf("Вставлено %d строк в таблицу 'tests_groups'.", affected)
				}
			}
		}
		defer CloseConnectToDB()
	}
	if err != nil {
		log.Errorf("Ошибка при добавлении новой Группы в БД: %v", err)
	}
	return err
}

// Удаляет из базы заданную группу
// Получает имя удаляемой группы
func DelTestGroup(groupName string) error {
	var err error
	var result sql.Result
	var affected int64
	SetLogFormat()

	// Проверить пермишен пользователя для удалений
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "delete_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Удаление Группы из базы
			log.Debugf("Удаление Группы: %s", groupName)
			result, err = db.Exec("DELETE FROM tests_groups WHERE name=?", groupName)
			if err == nil {
				affected, err = result.RowsAffected()
				if err == nil {
					if affected == 0 {
						err = errors.New(fmt.Sprintf("Ошибка удаления Группы '%s'. Есть такая Группа?", groupName))
						log.Errorf("Ошибка удаления Группы '%s'", groupName)
					}
					log.Debugf("Удалено строк в БД: %v", affected)
				}
			}
		}
		defer CloseConnectToDB()
	}
	if err != nil {
		log.Errorf("Ошибка при удалении Группы из БД: %v", err)
	}
	return err
}

// Изменяет имя группы тестов
// Получает имя редактируемой группы
func EditTestGroup(oldName string, newName string) error {

	var err error
	var result sql.Result
	var affected int64
	SetLogFormat()

	// Проверить пермишен пользователя для редактирования
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "edit_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Изменение имени Группы
			result, err = db.Exec("UPDATE tests_groups SET name=? WHERE name=?", newName, oldName)
			if err == nil {

				affected, err = result.RowsAffected()
				if err == nil {
					log.Debugf("Изменено строк: %v", affected)
				}

				if affected == 0 {
					err = errors.New(
						fmt.Sprintf("Ошибка изменения имени группы '%s' на новое имя '%s'", oldName, newName))
					log.Errorf("Ошибка изменения имени группы '%s' на новое имя '%s'", oldName, newName)
				}
			}
		}
		defer CloseConnectToDB()
	}
	if err != nil {
		log.Errorf("Ошибка при изменении имени Группы в БД: %v", err)
	}
	return err
}
