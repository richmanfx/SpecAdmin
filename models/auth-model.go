package models

type User struct {
	Login 		string
	Password	string
	FullName	string
	Permissions Permission
}

type Permission struct {
	Create 		bool		// Создавать Группы, Сюиты, Сценарии, Шаги
	Edit 		bool		// Редактировать Группы, Сюиты, Сценарии, Шаги
	Delete 		bool		// Удалять Группы, Сюиты, Сценарии, Шаги
	Config 		bool		// Конфигурировать SpecAdmin
	Users 		bool		// Создавать, удалять пользователей, редактировать их данные
}
