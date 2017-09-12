package helpers

import (
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"../../models"
	"database/sql"
	"github.com/jung-kurt/gofpdf"
	"fmt"
	"os"
	"strconv"
)

// Добавляет сценарий в БД
func AddTestScript(newScriptName string, scriptSerialNumber string, scriptSuite string) error {
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
			// Добавление Скрипта в БД
			log.Debugf("Добавление Скрипта: '%s', Порядковый номер '%s', Сюита: '%s'",newScriptName, scriptSerialNumber, scriptSuite)

			result, err = db.Exec("INSERT INTO tests_scripts (name, serial_number, name_suite) VALUES (?,?,?)", newScriptName, scriptSerialNumber, scriptSuite)
			if err == nil {
				affected, err = result.RowsAffected()
				if err == nil {
					log.Debugf("Вставлено %d строк в таблицу 'tests_scripts'.", affected)
				}
			}
		}
		db.Close()
	}
	if err != nil {log.Errorf("Ошибка при добавлении сценария '%s' в БД: '%v'", newScriptName, err)}
	return err
}


// Удаляет сценарий из БД
func DelTestScript(scriptName, scriptsSuiteName string) error {
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

			// Удаление скрипта из БД
			log.Debugf("Удаление Скрипта '%s'.", scriptName)
			result, err = db.Exec("DELETE FROM tests_scripts WHERE name=? AND name_suite=?", scriptName, scriptsSuiteName)
			if err == nil {
				var affected int64
				affected, err = result.RowsAffected()
				if err == nil {
					if affected == 0 {
						_, goModuleName, lineNumber, _ := runtime.Caller(1)
						log.Debugf("Ошибка удаления Скрипта '%s'. goModuleName=%v, lineNumber=%v",
							scriptName, goModuleName, lineNumber)
					}
					log.Debugf("Удалено строк в БД: %v.", affected)
				}
			}
		}
		db.Close()
	}
	if err != nil {log.Errorf("Ошибка при удалении сценария '%s': %v", scriptName, err)}
	return err
}

// Возвращает Сценарий из БД по имени Сценария и имени его Сюиты
func GetScript(scriptsName, scriptsSuiteName string) (models.Script, int, error) {
	var err error
	var script models.Script
	var id int
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Получить данные о Сценарии и его ключ 'id'
		log.Debugf("Получение данных Сценария '%s' из БД.", scriptsName)
		rows := db.QueryRow("SELECT id,serial_number FROM tests_scripts WHERE name=? AND name_suite=?",	scriptsName, scriptsSuiteName)

		// Получить данные из результата запроса

		err = rows.Scan(&id, &script.SerialNumber)
		if err == nil {
			log.Debugf("rows.Next из таблицы tests_scripts: %d, %s", id, script.SerialNumber)

			// Заполнить данными Сценарий
			script.Name = scriptsName
			script.Suite = scriptsSuiteName
		}
	}
	if err != nil {log.Errorf("Ошибка при получении данных Сценария '%s' Сюиты '%s' из БД: %v", scriptsName, scriptsSuiteName, err)}
	return script, id, err
}

// Обновить данные Сценария в БД
func UpdateTestScript(scriptId int, scriptName string, scriptSerialNumber string, scriptsSuite string) error {
	var err error
	SetLogFormat()

	// Проверить пермишен пользователя для редактирования
	log.Debugf("Проверка пермишена для пользователя '%s'", UserLogin)
	err = CheckOneUserPermission(UserLogin, "edit_permission")

	if err == nil {

		// Подключиться к БД
		err = dbConnect()
		if err == nil {
			// Обновить данные о Сценарии в БД
			log.Debugf("Обновление данных о Сценарии '%s' в БД", scriptName)
			_, err = db.Query("UPDATE tests_scripts SET name=?, serial_number=?, name_suite=? WHERE id=? LIMIT 1",
				scriptName, scriptSerialNumber, scriptsSuite, scriptId)
			if err == nil {
				log.Debugf("Успешно обновлены данные Сценария '%s' в БД.", scriptName)
			}
		}
	}
	if err != nil {log.Errorf("Ошибка обновления данных Сценария '%s' в БД: %v", scriptName, err)}
	return err
}


// Получить список имён Сюит в заданной Группе
func GetSuitsNameFromSpecifiedGroup(groupName string) ([]string, error) {
	var err error
	var name string
	var rows *sql.Rows
	suitsNameList := make([]string, 0, 0)
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Запрос имён Сюит
		rows, err = db.Query("SELECT name FROM tests_suits WHERE name_group=? ORDER BY serial_number", groupName)
		if err == nil {
			for rows.Next() {
				err = rows.Scan(&name)
				if err == nil {
					log.Debugf("rows.Next из таблицы tests_suits: %s", name)
					suitsNameList = append(suitsNameList, name) // Добавить в список
				}
			}
		}
	}
	db.Close()
	if err != nil {log.Errorf("Ошибка при получении списока имён Сюит в Группе '%s': %v", groupName, err)}
	return suitsNameList, err
}

// Получить только ID Сценариев для заданных Сюит
func GetScriptIdList(suitsNameFromGroup []string) ([]int, error) {

	var err error
	var id int
	var rows *sql.Rows
	scriptsIdList := make([]int, 0, 0)

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		for _, suiteName := range suitsNameFromGroup {
			rows, err = db.Query("SELECT id FROM tests_scripts WHERE name_suite=? ORDER BY serial_number", suiteName)
			if err == nil {

				// Данные получить из результата запроса
				for rows.Next() {
					err = rows.Scan(&id)
					if err == nil {
						log.Debugf("rows.Next из таблицы tests_scripts: %d", id)
						scriptsIdList = append(scriptsIdList, id)		// Добавить в список имён
					}
				}
			}
		}
	}
	db.Close()
	if err != nil {log.Errorf("Ошибка при получении ID Сценариев для Сюит: %v", err)}
	return scriptsIdList, err
}


// Получить Сценарии только для заданных Сюит
func GetScriptListForSpecifiedSuits(suitsNameFromGroup []string) ([]models.Script, error) {
	var err error
	var stepsList []models.Step
	var rows *sql.Rows
	scriptsList := make([]models.Script, 0, 0)

	// Получить только ID Сценариев для заданных Сюит
	scriptsIdList, err := GetScriptIdList(suitsNameFromGroup)
	if err == nil {

		// Получить Шаги из БД только для заданных по ID Сценариев
		stepsList, err = GetStepsListForSpecifiedScripts(scriptsIdList)
		if err == nil {

			// Подключиться к БД
			err = dbConnect()
			if err == nil {
				// Запрос Сценариев заданных Сюит из БД
				for _, suiteName := range suitsNameFromGroup {
					rows, err = db.Query(
						"SELECT id, name, serial_number, name_suite FROM tests_scripts WHERE name_suite=? ORDER BY serial_number", suiteName)
					if err == nil {
						// Данные получить из результата запроса
						for rows.Next() {
							var script models.Script
							err = rows.Scan(&script.Id, &script.Name, &script.SerialNumber, &script.Suite)
							if err == nil {

								log.Debugf("rows.Next из таблицы tests_scripts: %s, %s, %s, %s", script.Id, script.Name, script.SerialNumber, script.Suite)

								// Закинуть Шаги в Сценарий
								for _, step := range stepsList {
									if step.ScriptsId == script.Id { // Если Шаг принадлежит Сценарию, то добавляем его
										script.Steps = append(script.Steps, step)
										log.Debugf("Добавлен шаг '%v' в сценарий '%v'", step.Name, script.Name)
									} else {
										log.Debugf("Не добавлен шаг '%v' в сценарий '%v'", step.Name, script.Name)
									}
								}

								scriptsList = append(scriptsList, script) // Добавить сценарий в список
							}
							log.Debugf("Список сценариев: %v", scriptsList)
						}
					}
				}
			}
			db.Close()
		}
	}
	if err != nil {log.Errorf("Ошибка при получении Сценариев для заданных Сюит: %v", err)}
	return scriptsList, err
}


// Вернуть Сценарий и Сюиту Шага по ID Сценария
func GetScriptAndSuiteByScriptId(ScriptId int) (string, string, error) {
	var err error
	var stepsScriptName string
	var scripsSuiteName string
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Данные по Сценарию из БД
		log.Debugf("Получение данных из БД по Сценарию с Id='%s'.", ScriptId)
		err = db.QueryRow("SELECT name, name_suite FROM tests_scripts WHERE id=?", ScriptId).Scan(&stepsScriptName, &scripsSuiteName)

		if err == nil {
			log.Debugf("rows.Next из таблицы tests_scripts: %s, %s", stepsScriptName, scripsSuiteName)
		}
	}
	if err != nil {log.Errorf("Ошибка при получении данных из БД по Сценарию с Id='%s': %v", ScriptId, err)}
	return stepsScriptName, scripsSuiteName, err
}


// Сгенерировать PDF файл со списком Шагов Сценария
func GetScripsStepsPdf(scriptsSuite string, scriptName string, stepsList []models.Step) (string, error) {		// TODO: возвращать только ошибку без ПДФ-а

	var pdf *gofpdf.Fpdf
	var err error
	const pdfDirName string = "pdf"
	const pdfFileName string = "steps.pdf"

	pdf = gofpdf.New("L", "mm", "A4", "")

	pdf.SetFontLocation("fonts")
	fontSize := 12.0
	pdf.AddFont("Helvetica-1251", "B", "helvetica_1251.json")
	pdf.SetFont("Helvetica-1251", "B", fontSize)


	ht := pdf.PointConvert(fontSize)
	tr := pdf.UnicodeTranslatorFromDescriptor("cp1251")
	pageWidth, pageHeight := pdf.GetPageSize()
	leftMargin, rightMargin, _, bottomMargin := pdf.GetMargins()
	columnWidths := []float64{10, 100, 100, pageWidth - leftMargin - rightMargin - (10 + 100 + 100)}
	marginCell := 2. // margin of top/bottom of cell
	rows := [][]string{}

	write := func(str string) {
		pdf.MultiCell(pageWidth, ht, tr(str), "", "L", false)
		pdf.Ln(ht) 	// Переход на новую строку
	}
	pdf.AddPage()

	// Вывод контента в PDF-документ
	write(fmt.Sprintf("Сюита: %s", scriptsSuite))
	write(fmt.Sprintf("Сценарий: %s", scriptName))
	pdf.Ln(ht)

	// Заголовок таблицы
	header := []string{"N", "Название шага", "Описание шага", "Ожидаемый результат"}
	pdf.SetFillColor(200, 200, 200)
	for i, str := range header {
		pdf.CellFormat(columnWidths[i], ht+2, tr(str), "1", 0, "C", true, 0, "")
	}
	pdf.Ln(1.5 * ht)

	// Тело таблицы
	for _, step := range stepsList {
		rows = append(rows, []string{tr(strconv.Itoa(step.SerialNumber)), tr(step.Name), tr(step.Description), tr(step.ExpectedResult)})
	}
	splitAndWrapTextInCell(rows, pdf, columnWidths, marginCell, pageHeight, bottomMargin)

	// Записать файл
	fullPdfFileName := pdfDirName + string(os.PathSeparator) + pdfFileName
	err = pdf.OutputFileAndClose(fullPdfFileName)
	if err != nil {
		log.Errorf("Ошибка: '%v'", err)
	}

	return fullPdfFileName, err
}

// Разбивает текст для ячейки PDF-таблицы на слова, переносит слова и оборачивает в рамку ячейки.
func splitAndWrapTextInCell(rows [][]string, pdf *gofpdf.Fpdf, columnWidths []float64, marginCell float64, pageHeight float64, bottomMargin float64) {
	for _, row := range rows {
		curx, y := pdf.GetXY()
		x := curx

		height := 0.
		_, lineHt := pdf.GetFontSize()

		for i, txt := range row {
			lines := pdf.SplitLines([]byte(txt), columnWidths[i])
			h := float64(len(lines))*lineHt + marginCell*float64(len(lines))
			if h > height {
				height = h
			}
		}
		// Если высота строки не помещается на странице, то добавить новую страницу
		if pdf.GetY()+height > pageHeight-bottomMargin {
			pdf.AddPage()
			y = pdf.GetY()
		}
		for i, txt := range row {
			width := columnWidths[i]
			pdf.Rect(x, y, width, height, "")
			if i == 0 { // Нумерация шагов в центре ячейки
				pdf.MultiCell(width, lineHt+marginCell, txt, "", "C", false)
			} else {
				pdf.MultiCell(width, lineHt+marginCell, txt, "", "", false)
			}
			x += width
			pdf.SetXY(x, y)
		}
		pdf.SetXY(curx, y+height)
	}
}

// Получить ID Сценария по ID Шага
func GetScriptsIdByStepsId(stepsId int) (int, error) {
	var err error
	var scriptId int
	SetLogFormat()

	// Подключиться к БД
	err = dbConnect()
	if err == nil {

		// Данные по Сценарию из БД
		log.Debugf("Получение из БД ID Сценария по ID Шага '%s'.", stepsId)
		err = db.QueryRow("SELECT script_id FROM tests_steps WHERE id=?", stepsId).Scan(&scriptId)

		if err == nil {
			log.Debugf("rows.Next из таблицы tests_steps: %s", scriptId)
		}
	}
	if err != nil {log.Errorf("Ошибка при получении из БД шдентификатора Сценария по идентификатору Шага'%s': %v", stepsId, err)}
	return scriptId, err
}





















