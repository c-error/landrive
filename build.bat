@echo off

cd /d %~dp0

@REM :LOOP

echo.
echo Starting Build ...
echo.
echo - BUILDING: icon ...
windres icon.rc -O coff -o icon.syso

echo - BUILDING: drive.exe ...
go clean
go build -ldflags="-s -w" -trimpath

echo.
echo Cleaning ...
del /q icon.syso
echo.

echo Task Completed [Press any key to exit] ...
pause >nul 2>&1
exit
@REM goto LOOP

