@echo off
rem ��������� � ����������� �� BestServer ��� ������������ ������

rem ���������
rar a -R specadmin.rar specadmin css js templates sql images

rem ������� �� SSH
scp specadmin.rar zoer@bestserver:tmp/spec-admin/specadmin.rar

rem �������� ������
del specadmin.rar
