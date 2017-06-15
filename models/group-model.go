package models

type Step struct {
	Name 			string
	Number 			int			// Порядковый номер Шага
	Description 	string		// Что делать в Шаге
	ExpectedResult	string		// Ожидаемый результат
	ScreenShotPath	string		// Ссылка на скриншот
	Script			string		// Сценарий, которому принадлежит Шаг
}

type Script struct {
	Name 			string
	Number 			int			// Порядковый номер Сценария
	Suite 			string		// Сюита, которой принадлежит Сценарий
	Steps			[]Step		// Шаги
}

type Suite struct {
	Name			string
	Description 	string		// Описание Сюиты
	Group 			string		// Группа тестов, которой принадлежит Сюита
	Scripts 		[]Script	// Сценарии
}

type Group struct {
	Name		string
	Suites		[]Suite			// Сюиты
}
