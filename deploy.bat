﻿@echo off
rem Архивация и перемещение на BestServer для последующего деплоя

rem Архивация
rar a -R specadmin.rar specadmin css js templates

rem Заливка по SSH
scp specadmin.rar zoer@bestserver:tmp/spec-admin/specadmin.rar

rem Удаление архива
del specadmin.rar
