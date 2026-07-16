@echo off
setlocal enabledelayedexpansion
cd /d "%~dp0"

echo ===================================
echo   StreamFlix - Start
echo ===================================
echo.

where go >nul 2>nul
if errorlevel 1 (
    echo [FEHLER] Go wurde nicht gefunden. Bitte Go 1.24+ installieren:
    echo https://go.dev/dl/
    pause
    exit /b 1
)

if not exist ".env" (
    if exist ".env.example" (
        echo [INFO] Keine .env gefunden - kopiere .env.example nach .env
        copy /y ".env.example" ".env" >nul
        echo [HINWEIS] Bitte .env pruefen und Werte ^(DATABASE_URL, JWT_SECRET, ...^) anpassen.
        echo.
    ) else (
        echo [FEHLER] Weder .env noch .env.example gefunden.
        pause
        exit /b 1
    )
)

set /p RUN_SEED="Datenbank-Migrationen/Demo-Daten jetzt seeden? (j/N): "
if /i "%RUN_SEED%"=="j" (
    echo [INFO] Fuehre Seed aus...
    go run ./cmd/seed
    if errorlevel 1 (
        echo [FEHLER] Seed fehlgeschlagen. Ist PostgreSQL erreichbar und DATABASE_URL korrekt?
        pause
        exit /b 1
    )
    echo.
)

echo [INFO] Starte StreamFlix-Server...
echo.
go run ./cmd/streamflix

pause
