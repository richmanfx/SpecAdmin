package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

// Добавляет в базу новую сюиту тестов
// Получает имя новой сюиты
func AddTestSuite(suitesName string, suitesDescription string, suitesGroup string) error {

	var err error

	// Подключиться к БД
	err = dbConnect()
	if err != nil {
		return err
	}

	// Добавление Сюиты в базу, используем плейсхолдер
	log.Infof("Добавление Сюиты: %s, Описание: %s, Группа: %s", suitesName, suitesDescription, suitesGroup)
	result, err := db.Exec("INSERT INTO tests_suits (name, description, name_group) VALUES (?,?,?)", suitesName, suitesDescription, suitesGroup)

	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {panic(err)}
		log.Infof("Вставлено строк: %v", affected)
	}

	db.Close()

	return err
}