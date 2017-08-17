package helpers

import (
	"time"
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"../../models"
	"database/sql"
)

var UserLogin string

// Проверить наличие sessid в БД
func SessionIdExistInBD(sessidFromBrowser string) bool {

	var err error
	var sessidExist bool = false
	var result sql.Result
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Если expire у сессии истекло, то сессию удалить из БД и вернуть false
		var expires time.Time
		var user string

		requestResult := db.QueryRow("SELECT expires,user FROM sessions WHERE session_id=?", sessidFromBrowser)
		log.Debugf("requestResult: %v", requestResult)
		err = requestResult.Scan(&expires, &user)

		if err == nil {
			log.Debugf("'expires' из БД для sessid='%s': %v", sessidFromBrowser, expires)

			// Если expire позднее текущего времени, то сессия существует и валидна
			if expires.After(time.Now()) {
				log.Debugln("Сессия не истекла")
				sessidExist = true
				UserLogin = user // Имя пользователя в Хедер вывести
			} else {
				// Если сессия истекла, то удаляем её из БД
				result, err = db.Exec("DELETE FROM sessions WHERE session_id=?", sessidFromBrowser)
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
			log.Debugf("Сессии '%s' в таблице 'sessions' нет: %v", sessidFromBrowser, err)
			sessidExist = false
		}
	}
	if err != nil { log.Errorf("Ошибка при проверке наличия сессии в БД: %v", err) }
	return sessidExist
}


// Сохранить в БД Сессию
func SaveSessionInDB(sessid string, expires time.Time, userName string) error {

	var err error
	var affected int64
	var result sql.Result
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Внести в БД
		result, err = db.Exec("INSERT INTO sessions (session_id,expires,user) VALUE (?,?,?)", sessid, expires, userName)

		if err == nil {
			affected, err = result.RowsAffected()
			if err == nil {
				log.Debugf("Вставлено %d строк в таблицу 'sessions'.", affected)
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при сохранении Сессии в БД: '%v'", err)}
	return err
}


// Создать пользователя в БД
func CreateUserInDb(user models.User) error {

	var err error
	var affected int64
	var result sql.Result
	log.Debugf("user в 'CreateUserInDb': '%v'", user)

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Занести пользователя в БД
		result, err = db.Exec(
			"INSERT INTO user (login,passwd,salt,full_name,create_permission,edit_permission,delete_permission,config_permission,users_permission) VALUE (?,?,?,?,?,?,?,?,?)",
			user.Login, user.Password, user.Salt, user.FullName, user.Permissions.Create, user.Permissions.Edit,
			user.Permissions.Delete, user.Permissions.Config, user.Permissions.Users)

		if err == nil {
			affected, err = result.RowsAffected()
			if err == nil {
				log.Debugf("Вставлено %d строк в таблицу 'user'.", affected)
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при создании пользователя в БД: '%v'", err)}
	return err
}


// Сохранить пользователя в БД после редактирования
func SaveUserInDb(user models.User) error {

	var err error
	var affected int64
	var result sql.Result
	log.Debugf("user в 'SaveUserInDb': '%v'", user)

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	panic(err) }

	// Занести пользователя в БД
	result, err = db.Exec("UPDATE user SET full_name=?, create_permission=?, edit_permission=?, delete_permission=?, config_permission=?, users_permission=? WHERE login=? LIMIT 1",
		user.FullName, user.Permissions.Create, user.Permissions.Edit, user.Permissions.Delete, user.Permissions.Config, user.Permissions.Users, user.Login)

	if err == nil {
		affected, err = result.RowsAffected()
		if err == nil {
			log.Debugf("Изменено %d строк в таблице 'user'.", affected)
		}
	}

	if err != nil {log.Errorf("Ошибка при сохранении пользователя в БД после редактирования: '%v'", err)}
	return err
}


// Удалить пользователя из БД
func DeleteUserInDb(user models.User) error {

	var err error
	var result sql.Result
	log.Debugf("user в 'DeleteUserInDb': '%v'", user)

	// Проверить пермишен пользователя для удалений
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "delete_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Удалить пользователя в БД
			result, err = db.Exec("DELETE FROM user WHERE login=? AND full_name=?", user.Login, user.FullName)

			if err == nil {
				var affected int64
				affected, err = result.RowsAffected()

				if err == nil {
					if affected > 0 {
						log.Debugf("Удалено строк в БД: %d", affected)
					} else {
						err = fmt.Errorf("Ошибка удаления пользователя '%s' - не удалено ни одной строки в БД.", user.Login)
					}

				}
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при удалении пользователя '%s' в БД: '%v'", user.Login, err)}
	return err
}



// Считать из БД всех пользователей
func GetUsers() ([]models.User, error) {

	var err error
	var rows *sql.Rows
	usersList := make([]models.User, 0, 0)

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Считать
		rows, err = db.Query(
			"SELECT login, full_name, create_permission, edit_permission, delete_permission, config_permission, users_permission FROM user ORDER BY login")
		if err == nil {

			// Данные получить из результата запроса
			var user models.User
			for rows.Next() {
				err = rows.Scan(&user.Login, &user.FullName, &user.Permissions.Create, &user.Permissions.Edit,
					&user.Permissions.Delete, &user.Permissions.Config, &user.Permissions.Users)
				if err == nil {
					log.Debugf("User из таблицы user: %v", user)
					usersList = append(usersList, user) // Заполнить пользователями список пользователей
				}
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при считывании из БД всех пользователей: '%v'", err)}
	return usersList, err
}


// Проверка валидности пароля
func ValidatePassword(userName, oldPassword string) error {
	var err error
	var salt string
	var oldHash string

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Получить Соль из БД
		salt, err = GetSaltFromDb(userName)
		log.Debugf("Соль из БД: '%s'", salt)

		// Сгенерить Хеш пароля с Солью
		newHash := CreateHash(oldPassword, salt)
		log.Debugf("Хеш с Солью: '%s'", newHash)

		// Получить старый Хеш из БД
		oldHash, err = GetHashFromDb(userName)
		log.Debugf("Хеш из БД: '%s'", oldHash)

		if err == nil {
			// Сверить полученный Хеш с Хешем в БД
			if newHash != oldHash {
				log.Errorln("Хеш пароля не совпадает с хешем в БД")
				err = fmt.Errorf("Неверный пароль")
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при проверке пароля: '%v'", err)}
	return err
}


// Записать в БД новый пароль заданного пользователя
func SavePassword(userName, newPassword string) error {

	var err error
	var result sql.Result
	var affected int64

	// Получить Соль из БД
	salt, err := GetSaltFromDb(userName)
	log.Debugf("Соль из БД: '%s'", salt)

	// Сгенерить Хеш пароля с Солью
	newHash := CreateHash(newPassword, salt)
	log.Debugf("Хеш с Солью: '%s'", newHash)

	// Подключиться к БД
	err = dbConnect()
	if err != nil {

		// Занести новый хеш пароля в БД (соль не меняется)
		result, err = db.Exec("UPDATE user SET passwd=? WHERE login=?", newHash, userName)

		if err == nil {
			affected, err = result.RowsAffected()
			if err == nil {
				log.Debugf("Вставлено %d строк в таблицу 'user'.", affected)
			}
		}
	}
	if err != nil {log.Errorf("Ошибка записи в БД нового пароля пользователя: '%v'", err)}
	return err
}


// Проверить пароль по Хешу из БД
func CheckPasswordInDB(login, password string) error {

	// Получить Соль из БД
	salt, err := GetSaltFromDb(login)
	log.Debugf("Соль из БД: '%s'", salt)

	if err == nil {

		// Сгенерить Хеш пароля с Солью
		newHash := CreateHash(password, salt)
		log.Debugf("Хеш с Солью: '%s'", newHash)

		// Считать Хеш из БД
		var oldHash string
		oldHash, err = GetHashFromDb(login)

		if err == nil {
			if newHash == oldHash {
				log.Debugln("Хеш пароля совпадает с Хешем из БД.")
			} else {
				log.Errorln("Хеш пароля не совпадает с Хешем из БД.")
				err = fmt.Errorf("Неверный логин/пароль.")
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при проверке пароля по Хешу из БД: '%v'", err)}
	return err
}


// Проверить наличие пользователя в БД
func CheckUserInDB(login string) error {

	var err error
	var loginFromDB string

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Считать из БД
		requestResult := db.QueryRow("SELECT login FROM user WHERE login=?", login)

		err = requestResult.Scan(&loginFromDB)

		if err == nil {
			log.Debugf("Пользователь '%s' существует", login)
		} else {
			err = fmt.Errorf("Пользователь '%s' в БД не существует", login)
		}
	}
	if err != nil {log.Errorf("Ошибка при проверке наличия пользователя '%s' в БД: '%v'", login, err)}
	return err
}


// Удалить сессию из БД для заданного пользователя
func DeleteSession(userName string) error {

	var err error
	var result sql.Result

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Удалить сессию в БД
		result, err = db.Exec("DELETE FROM sessions WHERE user=?", userName)

		if err == nil {
			var affected int64
			affected, err = result.RowsAffected()

			if err == nil {
				if affected > 0 {
					log.Debugf("Удалено %d строк из таблицы 'sessions'", affected)
				} else {
					err = fmt.Errorf("Ошибка удаления сессии пользователя '%s' - не удалено ни одной строки из БД.", userName)
				}

			}
		}
	}
	if err != nil {log.Errorf("Ошибка удаления сессии пользователя '%s' из БД: '%v'", userName, err)}
	return err
}


// Проверить заданный пермишен у пользователя
func CheckOneUserPermission(login string, permission string) error {
	var err error
	var stmt *sql.Stmt

	// А существует ли пользователь?
	err = CheckUserInDB(login)

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Считать из БД
			requestString := fmt.Sprintf("SELECT %s FROM user WHERE login=?", permission)
			stmt, err = db.Prepare(requestString)
			requestResult := stmt.QueryRow(login)

			var permissionFromDB string		// Вовсе не bool !
			err = requestResult.Scan(&permissionFromDB)
			log.Debugf("permissionFromDB: '%v'", permissionFromDB)

			if err == nil {
				if permissionFromDB == "0" {
					err = fmt.Errorf("У пользователя '%s' недостаточно прав", login)
					log.Errorf("У пользователя '%s' недостаточно прав", login)
				} else {
					log.Debugf("У пользователя '%s' есть права '%s'", login, permission)
				}
			}
		}
	}
	if err != nil {log.Errorf("Ошибка проверки прав у пользователя '%s': '%v'", login, err)}
	return err
}



