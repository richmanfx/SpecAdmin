@echo off
rem Развернуть базу данных specadmin из дампа

mysql -p -u admin specadmin < specadmin.sql
