package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"../../models"
	"runtime"
	"database/sql"
	"github.com/jung-kurt/gofpdf"
	"os"
)

// Сформировать список Сюит для указанной Группы
func GetSuitesListInGroup(groupName string) ([]models.Suite, error) {
	var err error
	var rows *sql.Rows
	suitesList := make([]models.Suite, 0, 0)	// Слайс из Сюит
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Получить только список имён Сюит в данной Группе
		suitsNameFromGroup, err := GetSuitsNameFromSpecifiedGroup(groupName)
		if err == nil {
			log.Debugf("Сюиты из группы '%s': %v", groupName, suitsNameFromGroup)

			// Считать Сценарии только для заданных Сюит
			scriptsList, err := GetScriptListForSpecifiedSuits(suitsNameFromGroup)
			if err == nil {
				log.Debugf("Сценарии %v из Сюит %v", scriptsList, suitsNameFromGroup)

				// Формировать список Сюит
				for _, suiteName := range suitsNameFromGroup {

					// Сюиты из БД по списку имён Сюит
					rows, err = db.Query("SELECT name,description,serial_number FROM tests_suits WHERE name=? ORDER BY serial_number", suiteName)
					if err == nil {
						// Получить данные из результата запроса
						var suite models.Suite
						for rows.Next() {
							err = rows.Scan(&suite.Name, &suite.Description, &suite.SerialNumber)
							if err == nil {
								log.Debugf("Данные из таблицы 'tests_suits': %s, %s, %s, %s", suite.Name, suite.Description, suite.SerialNumber)

								// Заполнить Сюитами список Сюит
								suite.Group = groupName

								// Закинуть Сценарии в соответствующие Сюиты
								for _, script := range scriptsList { // Бежать по всем сценариям
									if script.Suite == suite.Name { // Если Сценарий принадлежит Сюите, то добавляем его
										suite.Scripts = append(suite.Scripts, script)
										log.Debugf("Добавлен сценарий '%v'('%v') в сюиту '%v'", script.Name, script.Suite, suite.Name)
									} else {
										log.Debugf("Не добавлен сценарий '%v'('%v') в сюиту '%v'", script.Name, script.Suite, suite.Name)
									}
								}

								// Добавить Сюиту в список
								suitesList = append(suitesList, suite)
							}
						}
					}
				}
			}
		}
	}
	log.Debugf("Список Сюит: '%v'", suitesList)
	if err != nil {log.Errorf("Ошибка формировании списка Сюит для Группы '%s': '%v'", groupName, err)}
	return suitesList, err
}


// Добавляет в базу новую сюиту тестов
// Получает имя новой сюиты, описание этой сюиты и группу тестов, в которую вносится сюита
func AddTestSuite(suitesName string, suitesDescription string, suitesSerialNumber string, suitesGroup string) error {
	var err error
	var result sql.Result
	var affected int64
	SetLogFormat()

	// Проверить пермишен пользователя для создания
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "create_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Добавление Сюиты в базу, используем плейсхолдер
			log.Debugf("Добавление Сюиты: %s, Описание: %s, Порядковый номер '%s' Группа: %s",
				suitesName, suitesDescription, suitesSerialNumber, suitesGroup)
			result, err = db.Exec("INSERT INTO tests_suits (name, description, serial_number, name_group) VALUES (?,?,?,?)",
				suitesName, suitesDescription, suitesSerialNumber, suitesGroup)
			if err == nil {
				affected, err = result.RowsAffected()
				if err == nil {
					log.Debugf("Вставлено %d строк в таблицу 'tests_suits'.", affected)
				}
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при добавлении новой Сюиты: '%v'", err)}
	return err
}

// Удаляет из базы сюиту
// Получает имя удаляемой сюиты
func DelTestSuite(suitesName string) error {
	var err error
	var result sql.Result
	SetLogFormat()

	// Проверить пермишен пользователя для удалений
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "delete_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {
			// Удаление Сюиты из базы данных
			log.Debugf("Удаление Сюиты '%s'.", suitesName)
			result, err = db.Exec("DELETE FROM tests_suits WHERE name=?", suitesName)
			if err == nil {
				var affected int64
				affected, err = result.RowsAffected()
				if err == nil {
					if affected == 0 {
						_, goModuleName, lineNumber, _ := runtime.Caller(1)
						err = fmt.Errorf("Ошибка удаления Сюиты '%s'. Есть такая Сюита?", suitesName)
						log.Debugf("Ошибка удаления Сюиты '%s'. goModuleName=%v, lineNumber=%v",
							suitesName, goModuleName, lineNumber)
					}
					log.Debugf("Удалено строк в БД: %v", affected)
				}
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при удалении из базы сюиты '%s': '%v'", suitesName, err)}
	return err
}

// Возвращает Сюиту из БД
func GetSuite(suitesName string) (models.Suite, error)  {
	var err error
	var suite models.Suite
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Получить данные о Сюите и её ключ 'id'
		log.Debugf("Получение данных Сюиты '%s' из БД.", suitesName)
		rows := db.QueryRow("SELECT serial_number, description, name_group FROM tests_suits WHERE name=?", suitesName)

		if err == nil {
			// Данные получить из результата запроса
			err = rows.Scan(&suite.SerialNumber, &suite.Description, &suite.Group)
			if err == nil {
				log.Debugf("rows.Next из таблицы tests_suits: %s, %s, %s", suite.Description, suite.SerialNumber, suite.Group)
				suite.Name = suitesName
			}
		}
	}
	if err != nil {log.Errorf("Ошибка при получении данных Сюиты '%s' из БД: '%v'", suitesName, err)}
	return suite, err
}

// Обновляет данные Сюиты в БД
func UpdateTestSuite(suitesName string, suitesDescription string,
					 suitesSerialNumber string, suitesGroup string) error {

	var err error
	var result sql.Result
	SetLogFormat()

	// Проверить пермишен пользователя для редактирования
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "edit_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {

			// Обновить данные о Сюите в БД
			log.Debugf("Обновление данных о Сюите '%s' в БД", suitesName)
			result, err = db.Exec("UPDATE tests_suits SET description=?, serial_number=?, name_group=? WHERE name=? LIMIT 1",
				suitesDescription, suitesSerialNumber, suitesGroup, suitesName)

			if err == nil {

				var affected int64
				affected, err = result.RowsAffected()
				if affected != 0 {
					log.Debugf("Успешно обновлены данные Сюиты '%s' в БД. Обновлено '%d' записей", suitesName, affected)
				} else {
					log.Errorf("Ошибка обновления данных Сюиты '%s' в БД. Обновлено '%d' записей", suitesName, affected)
					err = fmt.Errorf("Ошибка обновления данных Сюиты '%s' в БД. Обновлено '%d' записей", suitesName, affected)
				}

			}
		}
	}
	if err != nil { log.Errorf("Ошибка обновления данных Сюиты '%s' в БД: '%м'", suitesName, err) }
	return err
}


// Переименование Сюиты
func RenameTestSuite(oldSuiteName, newSuiteName string) error {

	var err error
	var result sql.Result
	SetLogFormat()

	// Проверить пермишен пользователя для редактирования
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "edit_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err != nil {

			// Обновить данные о Сюите в БД
			log.Debugf("Переименование Сюиты '%s' в '%s'", oldSuiteName, newSuiteName)
			result, err = db.Exec("UPDATE tests_suits SET name=? WHERE name=? LIMIT 1", newSuiteName, oldSuiteName)

			if err == nil {

				var affected int64
				affected, err = result.RowsAffected()
				if affected != 0 {
					log.Debugf("Успешно переименована Сюиты '%s' в '%s'. Обновлено '%d' записей", oldSuiteName, newSuiteName, affected)
				} else {
					log.Errorf("Ошибка переименования Сюиты '%s' в '%s'. Обновлено '%d' записей", oldSuiteName, newSuiteName, affected)
					err = fmt.Errorf("Ошибка переименования Сюиты '%s' в '%s'. Обновлено '%d' записей", oldSuiteName, newSuiteName, affected)
				}
			}
		}
	}
	if err != nil { log.Errorf("Ошибка переименования Сюиты '%s' в '%s': %v", oldSuiteName, newSuiteName, err) }
	return err
}

// Сгенерировать PDF файл со списком Сценариев Сюиты
func GetSuitesScriptsPdf(scriptsSuite string, scriptsList []models.Script) error {

	var pdf *gofpdf.Fpdf
	var err error
	const pdfDirName string = "pdf"
	const pdfFileName string = "scripts.pdf"

	pdf = gofpdf.New("L", "mm", "A4", "")

	pdf.SetFontLocation("fonts")
	fontSize := 12.0
	pdf.AddFont("Helvetica-1251", "B", "helvetica_1251.json")
	pdf.SetFont("Helvetica-1251", "B", fontSize)


	ht := pdf.PointConvert(fontSize)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	pageWidth, pageHeight := pdf.GetPageSize()
	leftMargin, rightMargin, _, bottomMargin := pdf.GetMargins()
	columnWidths := []float64{10, pageWidth - leftMargin - rightMargin - 10}
	marginCell := 2. // margin of top/bottom of cell
	rows := [][]string{}

	write := func(str string) {
		pdf.MultiCell(pageWidth, ht, tr(str), "", "L", false)
		pdf.Ln(ht) 	// Переход на новую строку
	}
	pdf.AddPage()

	// Вывод контента в PDF-документ
	write(fmt.Sprintf("Сюита: %s", scriptsSuite))
	pdf.Ln(ht)

	// Заголовок таблицы
	header := []string{"N", "Сценарий"}
	pdf.SetFillColor(200, 200, 200)
	for i, str := range header {
		pdf.CellFormat(columnWidths[i], ht+2, tr(str), "1", 0, "C", true, 0, "")
	}
	pdf.Ln(1.5 * ht)

	// Тело таблицы
	for _, script := range scriptsList {
		rows = append(rows, []string{tr(script.SerialNumber), tr(script.Name)})
	}
	splitAndWrapTextInCell(rows, pdf, columnWidths, marginCell, pageHeight, bottomMargin)

	// Записать файл
	fullPdfFileName := pdfDirName + string(os.PathSeparator) + pdfFileName
	err = pdf.OutputFileAndClose(fullPdfFileName)
	if err != nil {
		log.Errorf("Ошибка: '%v'", err)
	}

	return err
}

