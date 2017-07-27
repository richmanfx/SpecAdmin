package helpers

import (
	"time"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"../../models"
)


// Проверить наличие sessid в БД
func SessionIdExistInBD(sessidFromBrowser string) bool {

	var err error
	var sessidExist bool = false
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Если expire у сессии истекло, то сессию удалить из БД и вернуть false
	var expires time.Time

	requestResult := db.QueryRow("SELECT expires FROM sessions WHERE session_id=?", sessidFromBrowser)
	log.Infof("requestResult: %v", requestResult)
	err = requestResult.Scan(&expires)

	if err == nil {
		log.Infof("'expires' из БД для sessid='%s': %v", sessidFromBrowser, expires)

		// Если expire позднее текущего времени, то сессия существует и валидна
		if expires.After(time.Now()) {
			log.Infoln("Сессия не истекла")
			sessidExist = true
		} else {
			// Если сессия истекла, то удаляем её из БД
			result, err := db.Exec("DELETE FROM sessions WHERE session_id=?", sessidFromBrowser)
			if err == nil {
				var affected int64
				affected, err = result.RowsAffected()
				if err == nil {
					if affected == 0 {
						err = fmt.Errorf("Ошибка при удалениии сессии '%s'.", sessidFromBrowser)
						log.Debugf("Сессия '%s' НЕ удалена.", sessidFromBrowser)
					}
					log.Debugf("Удалено '%d' строк в таблице 'sessions'.", affected)
				}
			}
			sessidExist = false
		}
	} else {
		// Такой сессии в БД нет
		log.Infof("Сессии '%s' в таблице 'sessions' нет: %v", sessidFromBrowser, err)
		sessidExist = false
	}

	db.Close()

	return sessidExist
}


// Сохранить в БД Сессию
func SaveSessionInDB(sessid string, expires time.Time, userName string) error {

	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }


	// Внести в БД
	result, err := db.Exec("INSERT INTO sessions (session_id,expires,user) VALUE (?,?,?)", sessid, expires, userName)

	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {
			panic(err)
		}
		log.Debugf("Вставлено %d строк в таблицу 'sessions'.", affected)
	}

	db.Close()
	return err
}


// Создать пользователя в БД
func CreateUserInDb(user models.User) error {

	var err error
	log.Debugf("user в 'CreateUserInDb': '%v'", user)

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Занести пользователя в БД
	result, err := db.Exec("INSERT INTO user (login, passwd, full_name, create_permission, edit_permission, delete_permission, config_permission, users_permission) VALUE (?,?,?,?,?,?,?,?)", user.Login, user.Password, user.FullName, user.Permissions.Create, user.Permissions.Edit, user.Permissions.Delete, user.Permissions.Config, user.Permissions.Users)

	if err == nil {
		affected, err := result.RowsAffected()
		if err != nil {
			panic(err)
		}
		log.Debugf("Вставлено %d строк в таблицу 'user'.", affected)
	}

	db.Close()
	return err
}



// Удалить пользователя из БД
func DeleteUserInDb(user models.User) error {

	var err error
	log.Debugf("user в 'DeleteUserInDb': '%v'", user)

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Удалить пользователя в БД
	result, err := db.Exec("DELETE FROM user WHERE login=? AND full_name=?", user.Login, user.FullName)

	if err != nil {
		log.Errorf("Ошибка удаления пользователя '%s'", user.Login)
		return err
	} else {
		var affected int64
		affected, err = result.RowsAffected()

		if err == nil {
			if affected == 0 {
				log.Errorf("Ошибка удаления пользователя '%s'", user.Login)
				return err
			}
			log.Debugf("Удалено строк в БД: %d", affected)
		}
	}

	db.Close()
	return err
}



// Считать из БД всех пользователей
func GetUsers() ([]models.User, error) {

	var err error
	usersList := make([]models.User, 0, 0)

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Считать
	rows, err := db.Query("SELECT login, full_name, create_permission, edit_permission, delete_permission, config_permission, users_permission FROM user ORDER BY login")
	if err != nil {panic(err)}

	// Данные получить из результата запроса
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Login, &user.FullName, &user.Permissions.Create, &user.Permissions.Edit, &user.Permissions.Delete, &user.Permissions.Config, &user.Permissions.Users)
		if err != nil {panic(err)}
		log.Debugf("User из таблицы user: %v", user)

		// Заполнить пользователями список пользователей
		usersList = append(usersList, user)
	}

	db.Close()
	return usersList, err
}












