
Подготовка базы
===============

1. Установить MySQL

2. Создать новую базу 'specadmin':

    **create database specadmin;**

3. Создать нового пользователя:

    **create user 'specuser'@'localhost';**

4. Дать права пользователю на работу с базой:

    **grant select,insert,update,delete,index,alter,create,drop on specadmin.\* to specuser@localhost;**

5. Установить пароль новому пользователю:

    **ALTER USER 'specuser'@'localhost' IDENTIFIED BY 'Ghashiz7'; FLUSH PRIVILEGES;**

6. Залить пустую базу из дампа:

    **mysql -p -u some_admin specadmin < specadmin-without-data.sql**
    
    
Особенности
===========

1. Имена Сюит не могут содержать пробелы - id с пробелами не работают в JS.
    