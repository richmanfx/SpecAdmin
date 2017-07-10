package models

type Step struct {
	Id                 int
	Name               string
	SerialNumber       int    // Порядковый номер Шага
	Description        string // Что делать в Шаге
	ExpectedResult     string // Ожидаемый результат
	ScriptsId          int    // ID сценария, которому принадлежит Шаг
	ScreenShotFileName string // Имя скриншот-файла
}

type Script struct {
	Id 				int
	Name         	string
	SerialNumber 	string    	// Порядковый номер Сценария
	Suite        	string 		// Сюита, которой принадлежит Сценарий
	Steps        	[]Step 		// Шаги
}

type Suite struct {
	Name			string
	Description 	string		// Описание Сюиты
	SerialNumber	string		// Порядковый номер
	Group 			string		// Группа тестов, которой принадлежит Сюита
	Scripts 		[]Script	// Сценарии
}

type Group struct {
	Name		string
	Suites		[]Suite			// Сюиты
}
