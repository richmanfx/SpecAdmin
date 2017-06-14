
Подготовка базы
===============

1.Установить MySQL

2.Создать новую базу 'specadmin':
    create database specadmin;

3.Создать нового пользователя:
    create user 'specuser'@'localhost';

4.Дать права пользователю на работу с базой и задатьб пароль:
    grant select,insert,update,delete,index,alter,create,drop on specadmin.* to specuser@localhost identified by 'Ghashiz7';

5.Создать таблиц:
==========
create table tests_groups
(
	id int auto_increment
		primary key,
	name varchar(255) not null,
	constraint tests_groups_name_uindex
		unique (name)
)
comment 'Группы тестов'
;
==========
CREATE TABLE specadmin.tests_suits
(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    id_group INT NOT NULL,
    CONSTRAINT tests_suits_tests_groups__fk FOREIGN KEY (id_group) REFERENCES tests_groups (id)
);
CREATE UNIQUE INDEX tests_suits_name_uindex ON specadmin.tests_suits (name);
CREATE UNIQUE INDEX tests_suits_description_uindex ON specadmin.tests_suits (description);
ALTER TABLE specadmin.tests_suits COMMENT = 'Сюиты тестов';
==========