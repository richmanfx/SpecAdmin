package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
)

// Добавляет в базу новую сюиту тестов
// Получает имя новой сюиты
func AddTestSuite(suiteName string) error {

	var err error

	// Подключиться к БД
	err = dbConnect()
	if err != nil {
		return err
	}

	// Добавление Сюиты в базу, используем плейсхолдер
	//result, err := db.Exec("INSERT INTO tests_groups (name) VALUE (?)", suiteName)
	log.Infof("Сдесь будет добавление сюиты: %s", suiteName)
	//if err == nil {
	//	affected, err := result.RowsAffected()
	//	if err != nil {panic(err)}
	//	log.Infof("Вставлено строк: %v", affected)
	//}

	db.Close()

	return err
}