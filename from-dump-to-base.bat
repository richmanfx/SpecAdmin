@echo off
rem Восстановление базы "specadmin" из дампа

mysql -p -u admin specadmin < specadmin.sql
