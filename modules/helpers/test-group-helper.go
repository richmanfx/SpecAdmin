package helpers

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "fmt"
	"fmt"
)

var db *sql.DB

func GetTestGroupsList() []string {

	testGroupList := make([]string, 0, 30)
	var err error

	// Соединение с БД
	db, err = sql.Open("mysql", "specuser:Ghashiz7@tcp(localhost:3306)/specadmin?charset=utf8")
	if err != nil {panic(err)}

	// Проверка соединения с БД
	err = db.Ping()
	if err != nil {panic(err)}

	// Запрос Групп из базы, получить записи
	rows, err := db.Query("SELECT (name) FROM tests_groups")
	if err != nil {panic(err)}

	// Данные получить из результата запроса
	for rows.Next() {
		var group string
		err = rows.Scan(&group)
		if err != nil {panic(err)}
		//fmt.Println("row.Next group: ", group)
		testGroupList = append(testGroupList, group)
	}

	rows.Close()
	db.Close()

	return testGroupList
}

// Добавляет в базу новую группу тестов
// Получает имя новой группы
func AddTestGroup(groupName string) error {

	var err error

	// Соединение с БД
	db, err = sql.Open("mysql", "specuser:Ghashiz7@tcp(localhost:3306)/specadmin?charset=utf8")
	if err != nil {panic(err)}

	// Проверка соединения с БД
	err = db.Ping()
	if err != nil {panic(err)}

	// Добавление Группы в базу, используем плейсхолдер
	result, err := db.Exec("INSERT INTO tests_groups (name) VALUE (?)", groupName)
	if err != nil {panic(err)}

	affected, err := result.RowsAffected()
	if err != nil {panic(err)}
	fmt.Println("Вставлено строк: ", affected)

	db.Close()

	return err
}


