
Подготовка базы
===============

1. Установить MySQL

2. Создать новую базу 'specadmin':

    **create database specadmin;**

3. Создать нового пользователя:
    create user 'specuser'@'localhost';

4. Дать права пользователю на работу с базой и задать пароль:

    **grant select,insert,update,delete,index,alter,create,drop on specadmin.* to specuser@localhost identified by 'super_password';**

5. Залить пустую базу из дампа:

    **mysql -p -u some_admin specadmin < specadmin-without-data.sql**