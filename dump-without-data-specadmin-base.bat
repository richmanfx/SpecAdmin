@echo off
rem Дамп структуры базы данных specadmin без данных

mysqldump -u admin -p --routines --no-data --no-create-db --skip-opt specadmin > specadmin-without-data.sql
