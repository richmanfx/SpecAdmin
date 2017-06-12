
Подготовка базы
===============

1.Установить MySQL

2.Создать новую базу 'specadmin':
    create database specadmin;

3.Создать нового пользователя:
    create user 'specuser'@'localhost';

4.Дать права пользователю на работу с базой и задатьб пароль:
    grant select,insert,update,delete,index,alter,create,drop on specadmin.* to specuser@localhost identified by 'Ghashiz7';

5.Создать таблицы:

CREATE TABLE specadmin.tests_groups
			(
			id INT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(255) NOT NULL
			)

ALTER TABLE specadmin.tests_groups COMMENT = 'Группы тестов'

