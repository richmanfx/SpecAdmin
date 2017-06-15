@echo off
rem Дамп базы данных specadmin

mysqldump -p -u admin specadmin > specadmin.sql
