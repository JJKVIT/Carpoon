@echo off
set "BATDIR=%~dp0"

"%BATDIR%..\carpoon.exe" %* 2> "%temp%\carpoon_path.txt"

set "last="
for /f "usebackq delims=" %%d in ("%temp%\carpoon_path.txt") do set "last=%%d"

del "%temp%\carpoon_path.txt"


if defined last (
    cd /d "%last%"
    cls
)
