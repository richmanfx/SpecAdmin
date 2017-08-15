package helpers

import (
	"../../models"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"fmt"
	"os"
	"database/sql"
)

// Получить из БД конфигурационные данные
func GetConfig() ([]models.Option, error) {
	var err error
	var rows *sql.Rows
	var option models.Option
	config := make([]models.Option, 0, 0)
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Получить из БД конфигурационные параметры
		log.Debugln("Получение конфигурационных параметров из БД.")
		rows, err = db.Query("SELECT * FROM configuration ORDER BY id")

		if err == nil {
			// Данные получить из результата запроса

			for rows.Next() {
				err = rows.Scan(&option.Id, &option.Name, &option.Value)
				if err == nil {
					log.Debugf("rows.Next из таблицы configuration: %d %s %s", option.Id, option.Name, option.Value)
					config = append(config, option)
				}
			}
		}
	}
	db.Close()
	if err != nil {log.Errorf("Ошибка при получении конфигурационных данных из БД: '%v'", err)}
	return config, err
}


// Сохранить конфигурацию в БД
func SaveConfig(screenShotPath string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Записать в БД
		log.Debugln("Запись в БД пути к директории со скриншотами.")
		_, err = db.Query("UPDATE configuration SET property_value=? WHERE property_name=? LIMIT 1", screenShotPath, "Путь к скриншотам")

		// Здесь будет запись других параметров
		// ************************************

		if err == nil {
			log.Debugln("Конфигурационные параметры успешно сохранены в БД.")
		}
	}
	db.Close()
	if err != nil {log.Errorf("Ошибка при сохранении конфигурационных параметров в БД: '%v'", err)}
	return err
}


// Получить Путь к скриншотам
func GetScreenShotsPath() (string, error) {

	var screenShotsPath string
	var err error
	var config []models.Option

	config, err = GetConfig() // Получить из базы все конфигурационные данные
	if err == nil {
		for _, configItem := range config { // Выбрать про путь к скриншотам
			if configItem.Name == "Путь к скриншотам" {
				screenShotsPath = configItem.Value
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при получении пути к скриншотам: '%v'", err)}
	return screenShotsPath, err
}


// Получить список имён неиспользуемых файлов скриншотов
func GetUnusedScreenShotsFileName() ([]string, error) {

	unusedScreenShotsFileNameList := make([]string, 0, 0)

	// Получить все используемые имена файлов скриншотов
	usedScreenShotsFileNameList, err := GetScreenShotsFileName()
	if err == nil {
		// Получить список имён файлов из хранилища
		var screenShotsFileNameInRepositoryList []string
		screenShotsFileNameInRepositoryList, err = GetScreenShotsFileNameInRepositoryList()
		if err == nil {
			// Скрыжить имена файлов из хранилища и из БД
			for _, inRepoFileName := range screenShotsFileNameInRepositoryList {
				var matchesFlag bool = false
				for _, usedFileName := range usedScreenShotsFileNameList {
					if inRepoFileName == usedFileName {
						matchesFlag = true
						break
					}
				}
				if matchesFlag == false {
					unusedScreenShotsFileNameList = append(unusedScreenShotsFileNameList, inRepoFileName)
				}
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при получении списка имён файлов скриншотов в хранилище: '%v'", err)}
	return  unusedScreenShotsFileNameList, err
}


// Получить все используемые имена файлов скриншотов
func GetScreenShotsFileName() ([]string, error) {
	var err error
	var rows *sql.Rows
	usedScreenShotsFileNameList := make([]string, 0, 0)
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {
		rows, err = db.Query("SELECT screen_shot_file_name FROM tests_steps ORDER BY screen_shot_file_name")

		if err == nil {
			// Данные получить из результата запроса
			var usedScreenShotsFileName string
			for rows.Next() {
				err = rows.Scan(&usedScreenShotsFileName)
				if err == nil {
					log.Debugf("rows.Next из таблицы tests_steps: %s", usedScreenShotsFileName)

					// Дополнить список файлов
					usedScreenShotsFileNameList = append(usedScreenShotsFileNameList, usedScreenShotsFileName)
				}
			}
		}
	}
	db.Close()
	if err != nil {log.Errorf("Ошибка при получении всех используемыех имён файлов скриншотов: '%v'", err)}
	return usedScreenShotsFileNameList, err
}


// Получить список имён файлов из хранилища
func GetScreenShotsFileNameInRepositoryList() ([]string, error) {

	var err error
	var screenShotsPath string
	var files []os.FileInfo
	inRepositoryScreenShotsFileNameList := make([]string, 0, 0)
	SetLogFormat()

	// Получить Путь к скриншотам
	screenShotsPath, err = GetScreenShotsPath()
	if err == nil {
		// Получить список файлов
		files, err = ioutil.ReadDir(screenShotsPath)
		if err == nil {
			for _, file := range files {
				fmt.Println(file.Name())
				inRepositoryScreenShotsFileNameList = append(inRepositoryScreenShotsFileNameList, file.Name())
			}
			log.Debugf("Получен из директории '%s' список имен файлов: '%v' ", screenShotsPath, inRepositoryScreenShotsFileNameList)
		}
	}
	if err != nil {log.Errorf("Ошибка при получении списка имён файлов из директории '%s': '%v'", screenShotsPath, err)}
	return inRepositoryScreenShotsFileNameList, err
}
