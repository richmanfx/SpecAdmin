package helpers

import (
	"time"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"../../models"
)

var UserLogin string

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
	var user string

	requestResult := db.QueryRow("SELECT expires,user FROM sessions WHERE session_id=?", sessidFromBrowser)
	log.Infof("requestResult: %v", requestResult)
	err = requestResult.Scan(&expires, &user)

	if err == nil {
		log.Infof("'expires' из БД для sessid='%s': %v", sessidFromBrowser, expires)

		// Если expire позднее текущего времени, то сессия существует и валидна
		if expires.After(time.Now()) {
			log.Infoln("Сессия не истекла")
			sessidExist = true
			UserLogin = user		// Имя пользователя в Хедер вывести
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
		log.Infof("Вставлено %d строк в таблицу 'sessions'.", affected)
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
	result, err := db.Exec(
		"INSERT INTO user (login, passwd, salt, full_name, create_permission, edit_permission, delete_permission, " +
			"config_permission, users_permission) VALUE (?,?,?,?,?,?,?,?,?)",
		user.Login, user.Password, user.Salt, user.FullName, user.Permissions.Create, user.Permissions.Edit,
		user.Permissions.Delete, user.Permissions.Config, user.Permissions.Users)

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


// Сохранить пользователя в БД после редактирования
func SaveUserInDb(user models.User) error {

	var err error
	log.Debugf("user в 'SaveUserInDb': '%v'", user)

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Занести пользователя в БД
	result, err := db.Exec("UPDATE user SET full_name=?, create_permission=?, edit_permission=?, delete_permission=?, config_permission=?, users_permission=? WHERE login=? LIMIT 1",
		user.FullName, user.Permissions.Create, user.Permissions.Edit, user.Permissions.Delete, user.Permissions.Config, user.Permissions.Users, user.Login)

	if err == nil {
		affected, err := result.RowsAffected()
		if err == nil {
			log.Debugf("Изменено %d строк в таблице 'user'.", affected)
		}
	}

	db.Close()
	return err
}


// Удалить пользователя из БД
func DeleteUserInDb(user models.User) error {

	var err error
	log.Debugf("user в 'DeleteUserInDb': '%v'", user)

	// Проверить пермишен пользователя для удалений
	log.Infof("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "delete_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err != nil {
			panic(err)
		}

		// Удалить пользователя в БД
		result, _ := db.Exec("") 	// Накостылил здесь
		result, err = db.Exec("DELETE FROM user WHERE login=? AND full_name=?", user.Login, user.FullName)

		if err == nil {
			var affected int64
			affected, err = result.RowsAffected()

			if err == nil {
				if affected == 0 {
					log.Errorf("Ошибка удаления пользователя '%s'", user.Login)
					return err
				}
				log.Debugf("Удалено строк в БД: %d", affected)
			}
		} else {
			log.Errorf("Ошибка удаления пользователя '%s'", user.Login)
		}

		db.Close()
	}
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
	rows, err := db.Query(
		"SELECT login, full_name, create_permission, edit_permission, delete_permission, config_permission, users_permission FROM user ORDER BY login")
	if err != nil {panic(err)}

	// Данные получить из результата запроса
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.Login, &user.FullName, &user.Permissions.Create, &user.Permissions.Edit,
						&user.Permissions.Delete, &user.Permissions.Config, &user.Permissions.Users)
		if err != nil {panic(err)}
		log.Debugf("User из таблицы user: %v", user)

		// Заполнить пользователями список пользователей
		usersList = append(usersList, user)
	}

	db.Close()
	return usersList, err
}


// Проверка валидности пароля
func ValidatePassword(userName, oldPassword string) error {
	var err error
	var salt string

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }


	// Получить Соль из БД
	salt, err = GetSaltFromDb(userName)
	log.Infof("Соль из БД: '%s'", salt)


	// Сгенерить Хеш пароля с Солью
	newHash := CreateHash(oldPassword, salt)
	log.Infof("Хеш с Солью: '%s'", newHash)

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Получить старый Хеш из БД
	oldHash, err := GetHashFromDb(userName)
	log.Infof("Хеш из БД: '%s'", oldHash)

	if err == nil {
		// Сверить полученный Хеш с Хешем в БД
		if newHash != oldHash {
			log.Errorln("Хеш пароля не совпадает с хешем в БД")
			err = fmt.Errorf("Неверный пароль")
		}
	}

	db.Close()
	return err
}


// Записать в БД новый пароль заданного пользователя
func SavePassword(userName, newPassword string) error {

	// Получить Соль из БД
	salt, err := GetSaltFromDb(userName)
	log.Infof("Соль из БД: '%s'", salt)

	// Сгенерить Хеш пароля с Солью
	newHash := CreateHash(newPassword, salt)
	log.Infof("Хеш с Солью: '%s'", newHash)

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Занести новый хеш пароля и соль в БД
	result, err := db.Exec("UPDATE user SET passwd=?,salt=? WHERE login=?", newHash, salt, userName)

	if err == nil {
		affected, err := result.RowsAffected()
		if err == nil {
			log.Debugf("Вставлено %d строк в таблицу 'user'.", affected)
		}
	}
	db.Close()
	return err
}


// Проверить пароль по Хешу в БД
func CheckPasswordInDB(login, password string) error {

	// Получить Соль из БД
	salt, err := GetSaltFromDb(login)
	log.Infof("Соль из БД: '%s'", salt)

	if err == nil {

		// Сгенерить Хеш пароля с Солью
		newHash := CreateHash(password, salt)
		log.Infof("Хеш с Солью: '%s'", newHash)

		// Считать Хеш из БД
		var oldHash string
		oldHash, err = GetHashFromDb(login)

		if err == nil {
			if newHash == oldHash {
				log.Infoln("Хеш пароля совпадает с Хешем из БД.")
			} else {
				log.Errorln("Хеш пароля не совпадает с Хешем из БД.")
				err = fmt.Errorf("Не верный логин/пароль.")
			}
		}
	}
	return err
}


// Проверить наличие пользователя в БД
func CheckUserInDB(login string) error {

	var loginFromDB string

	// Подключиться к БД
	err := dbConnect()
	if err != nil {	panic(err) }

	// Считать из БД
	requestResult := db.QueryRow("SELECT login FROM user WHERE login=?", login)

	err = requestResult.Scan(&loginFromDB)

	if err == nil {
		log.Infof("Пользователь '%s' существует", login)
	} else {
		log.Infof("Пользователь '%s' в БД не существует", login)
		err = fmt.Errorf("Пользователь '%s' в БД не существует", login)
	}

	db.Close()
	return err
}


// Удалить сессию из БД для заданного пользователя
func DeleteSession(userName string) error {

	// Подключиться к БД
	err := dbConnect()
	if err != nil {	panic(err) }

	// Удалить сессию в БД
	result, err := db.Exec("DELETE FROM sessions WHERE user=?", userName)

	if err != nil {
		log.Errorf("Ошибка удаления сессии пользователя '%s'", userName)
		return err
	} else {
		var affected int64
		affected, err = result.RowsAffected()

		if err == nil {
			if affected == 0 {
				log.Errorf("Ошибка удаления сессии пользователя '%s'", userName)
				return err
			}
			log.Infof("Удалено %d строк из таблицы 'sessions'", affected)
		}
	}

	db.Close()
	return err
}


// Проверить заданный пермишен у пользователя
func CheckOneUserPermission(login string, permission string) error {

	var err error


	// А существует ли пользователь?
	err = CheckUserInDB(login)

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Считать из БД
			stmt, _ := db.Prepare("")	// По другому не с мог определить переменную :-)
			requestString := fmt.Sprintf("SELECT %s FROM user WHERE login=?", permission)
			stmt, err = db.Prepare(requestString)
			requestResult := stmt.QueryRow(login)

			var permissionFromDB string		// Вовсе не bool !
			err = requestResult.Scan(&permissionFromDB)
			log.Infof("permissionFromDB: '%v'", permissionFromDB)

			if err == nil {

				if permissionFromDB == "0" {
					err = fmt.Errorf("У пользователя '%s' недостаточно прав.", login)
					log.Infof("У пользователя '%s' недостаточно прав.", login)
				} else {
					log.Infof("У пользователя '%s' есть права '%s'.", login, permission)
				}
			}
		}
	}
	db.Close()
	return err
}



