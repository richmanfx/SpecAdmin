@echo off
rem ���� ��������� ���� ������ specadmin ��� ������

mysqldump -u admin -p --routines --no-data --no-create-db --skip-opt specadmin > specadmin-without-data.sql
