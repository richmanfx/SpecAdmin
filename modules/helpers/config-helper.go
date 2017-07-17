package helpers

import (
	"../../models"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"fmt"
)

// Получить из БД конфигурационные данные
func GetConfig() ([]models.Option, error) {
	var err error
	config := make([]models.Option, 0, 0)
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return config, err }

	// Получить из БД конфигурационные параметры
	log.Debugln("Получение конфигурационных параметров из БД.")
	rows, err := db.Query("SELECT * FROM configuration ORDER BY id")

	if err == nil {
		// Данные получить из результата запроса
		var propertyId int
		var propertyName string
		var propertyValue string

		for rows.Next() {
			err = rows.Scan(&propertyId, &propertyName, &propertyValue)
			log.Debugf("rows.Next из таблицы configuration: %d %s %s", propertyId, propertyName, propertyValue)

			// Заполнить параметрами конфигурационный список
			var option models.Option
			option.Id = propertyId
			option.Name = propertyName
			option.Value = propertyValue
			config = append(config, option)
		}
	}
	db.Close()
	return config, err
}


// Сохранить конфигурацию в БД
func SaveConfig(screenShotPath string) error {
	var err error
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return err }

	// Записать в БД
	log.Infoln("Запись в БД пути к директории со скриншотами.")
	_, err = db.Query("UPDATE configuration SET property_value=? WHERE property_name=? LIMIT 1", screenShotPath, "Путь к скриншотам")

	// Здесь будет запись других параметров
	// ************************************

	if err == nil {
		log.Debugln("Конфигурационные параметры успешно сохранены в БД.")
	} else {
		log.Debugln("Ошибка при сохранении конфигурационных параметров в БД.")
	}

	db.Close()
	return err
}


// Получить Путь к скриншотам
func GetScreenShotsPath() string {

	var screenShotsPath string

	config, err := GetConfig() // Получить из базы все конфигурационные данные
	if err != nil {
		panic(err)
	}
	for _, configItem := range config { // Выбрать про путь к скриншотам
		if configItem.Name == "Путь к скриншотам" {
			screenShotsPath = configItem.Value
		}
	}
	return screenShotsPath
}


// Получить список имён неиспользуемых файлов скриншотов
func GetUnusedScreenShotsFileName() ([]string, error) {

	unusedScreenShotsFileNameList := make([]string, 0, 0)

	// Получить все используемые имена файлов скриншотов
	usedScreenShotsFileNameList, err := GetScreenShotsFileName()
	if err != nil {
		log.Errorln("Ошибка при получении списка используемых имён файлов скриншотов из БД.")
	} else {

		// Получить список имён файлов из хранилища
		screenShotsFileNameInRepositoryList, err := GetScreenShotsFileNameInRepositoryList()
		if err != nil {
			log.Errorln("Ошибка при получении списка имён файлов скриншотов в хранилище.")
		} else {
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
	return  unusedScreenShotsFileNameList, err
}


// Получить все используемые имена файлов скриншотов
func GetScreenShotsFileName() ([]string, error) {
	var err error
	usedScreenShotsFileNameList := make([]string, 0, 0)
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err != nil {	return usedScreenShotsFileNameList, err }

	rows, err := db.Query("SELECT screen_shot_file_name FROM tests_steps ORDER BY screen_shot_file_name")

	if err == nil {
		// Данные получить из результата запроса
		var usedScreenShotsFileName string
		for rows.Next() {
			err = rows.Scan(&usedScreenShotsFileName)
			log.Infof("rows.Next из таблицы tests_steps: %s", usedScreenShotsFileName)

			// Дополнить список файлов
			usedScreenShotsFileNameList = append(usedScreenShotsFileNameList, usedScreenShotsFileName)
		}
	}
	db.Close()
	return usedScreenShotsFileNameList, err
}


// Получить список имён файлов из хранилища
func GetScreenShotsFileNameInRepositoryList() ([]string, error) {

	var err error
	inRepositoryScreenShotsFileNameList := make([]string, 0, 0)
	SetLogFormat()

	// Получить Путь к скриншотам
	screenShotsPath := GetScreenShotsPath()

	// Получить список файлов
	files, err := ioutil.ReadDir(screenShotsPath)
	if err != nil {
		log.Errorf("Ошибка при получении списка имён файлов из директории '%s'.", screenShotsPath)
	}

	for _, file := range files {
		fmt.Println(file.Name())
		inRepositoryScreenShotsFileNameList = append(inRepositoryScreenShotsFileNameList, file.Name())
	}

	log.Infof("Получен из директории '%s' список имен файлов: '%v' ", screenShotsPath, inRepositoryScreenShotsFileNameList)
	return inRepositoryScreenShotsFileNameList, err
}
