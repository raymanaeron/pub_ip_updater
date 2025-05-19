@echo off
echo Building Public IP Updater for DigitalOcean DNS...

REM Check if build directory exists, create if not
if not exist build mkdir build

REM Build the executable
go build -o build\pub_ip_updater.exe main.go

REM Check if build was successful
if %ERRORLEVEL% EQU 0 (
    echo.
    echo Build successful! Executable created at build\pub_ip_updater.exe
) else (
    echo.
    echo Build failed with error code %ERRORLEVEL%.
)
