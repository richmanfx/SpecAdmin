package models

import "google.golang.org/api/identitytoolkit/v3"

type Step struct {
	Id 			   int
	Name           string
	SerialNumber   int			// Порядковый номер Шага
	Description    string		// Что делать в Шаге
	ExpectedResult string		// Ожидаемый результат
	ScreenShotPath string		// Ссылка на скриншот
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
