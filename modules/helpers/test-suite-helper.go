package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"../../models"
	"runtime"
)



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
		log.Debugf("Вставлено строк: %v", affected)
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
	log.Debugf("Удаление Сюиты: %s", suitesName)
	result, err := db.Exec("DELETE FROM tests_suits WHERE name=?", suitesName)
	if err == nil {
		var affected int64
		affected, err = result.RowsAffected()
		if err == nil {
			if affected == 0 {
				_, fn, line, _ := runtime.Caller(1)
				err = fmt.Errorf("Ошибка удаления Сюиты '%s'. Есть такая Сюита?", suitesName)
				log.Debugf("Ошибка удаления Сюиты '%s'. file=%v, line=%v", suitesName, fn, line)
			}
			log.Debugf("Удалено строк в БД: %v", affected)
		}
	}
	db.Close()
	return err
}

// Получает Сюиту из БД
func GetSuite(suitesName string) (models.Suite, int, error)  {

	var err error
	var suite models.Suite

	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return suite, 0, err }

	// Получить данные о Сюите и её ключ 'id'
	log.Debugf("Получение данных Сюиты '%s' из БД", suitesName)
	rows, err := db.Query("SELECT id,description, serial_number ,name_group FROM tests_suits WHERE name=?", suitesName)
	if err != nil {panic(err)}			// TODO: Обработку сделать и вывод в браузер

	// Данные получить из результата запроса
	var id int
	var description string
	var serial_number string
	var name_group string
	for rows.Next() {
			err = rows.Scan(&id, &description, &serial_number, &name_group)
			if err != nil {panic(err)}
			log.Debugf("rows.Next из таблицы tests_suits: %d, %s, %s, %s", id, description, serial_number, name_group)
	}

	// Заполнить данными Сюиту
	suite.Name = suitesName
	suite.Description = description
	suite.SerialNumber = serial_number
	suite.Group = name_group

		db.Close()
	return suite, id, err
}

// Обновляет данные Сюиты в БД
func UpdateTestSuite(suitesId int, suitesName string, suitesDescription string,
					 suitesSerialNumber string, suitesGroup string) error {

	var err error

	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {
		return err
	}

	// Обновить данный о Сюите в БД
	log.Debugf("Обновление данных о Сюите '%s' в БД", suitesName)
	_, err = db.Query("UPDATE tests_suits SET name=?, description=?, serial_number=?, name_group=? WHERE id=? LIMIT 1",
		suitesName, suitesDescription, suitesSerialNumber, suitesGroup, suitesId)
	if err == nil {
		log.Debugf("Успешно обновлены данные Сюиты '%s' в БД.", suitesName)
	} else {
		log.Debugf("Ошибка обновления данныех Сюиты '%s' в БД.", suitesName)
	}

	return err
}
