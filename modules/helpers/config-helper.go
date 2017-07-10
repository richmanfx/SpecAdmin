package helpers

import (
	"../../models"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
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
	log.Infoln("Получение конфигурационных параметров из БД.")
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