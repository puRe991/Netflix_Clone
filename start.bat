@echo off
setlocal enabledelayedexpansion
cd /d "%~dp0"

echo ===================================
echo   StreamFlix - Start
echo ===================================
echo.

rem --- winget verfuegbar? (wird fuer Auto-Installation benoetigt) ---
where winget >nul 2>nul
if errorlevel 1 (
    set HAS_WINGET=0
) else (
    set HAS_WINGET=1
)

rem --- Echtes 32-Bit-Windows erkennen (nicht nur eine 32-Bit-cmd auf 64-Bit-Windows) ---
set IS_32BIT=0
if /i "%PROCESSOR_ARCHITECTURE%"=="x86" if not defined PROCESSOR_ARCHITEW6432 set IS_32BIT=1

rem --- Go pruefen, bei Bedarf installieren ---
where go >nul 2>nul
if errorlevel 1 (
    echo [INFO] Go wurde nicht gefunden.
    if "!HAS_WINGET!"=="1" (
        echo [INFO] Installiere Go ueber winget ...
        winget install --id GoLang.Go -e --accept-source-agreements --accept-package-agreements
        if !ERRORLEVEL! neq 0 (
            echo [FEHLER] Go-Installation ueber winget fehlgeschlagen.
            echo Bitte manuell installieren: https://go.dev/dl/
            pause
            exit /b 1
        )
        echo [HINWEIS] Go wurde installiert. Bitte dieses Fenster schliessen und
        echo start.bat erneut ausfuehren, damit die PATH-Aenderung wirksam wird.
        pause
        exit /b 0
    ) else (
        echo [FEHLER] winget nicht verfuegbar - Go kann nicht automatisch installiert werden.
        echo Bitte manuell installieren: https://go.dev/dl/
        pause
        exit /b 1
    )
) else (
    for /f "tokens=3" %%v in ('go version') do echo [OK] Go gefunden: %%v
)

rem --- PostgreSQL-Client (psql) pruefen, bei Bedarf installieren ---
where psql >nul 2>nul
if errorlevel 1 (
    echo [INFO] PostgreSQL-Client ^(psql^) wurde nicht gefunden.
    if "!IS_32BIT!"=="1" (
        rem PostgreSQL liefert seit Jahren keine 32-Bit-Windows-Installer mehr
        rem (auch nicht ueber winget) - ein lokaler Auto-Install kann auf einem
        rem echten 32-Bit-System also grundsaetzlich nicht funktionieren.
        echo [HINWEIS] Dieses System ist 32-Bit. PostgreSQL bietet dafuer keine
        echo Installer mehr an ^(auch winget findet keinen passenden^) - eine lokale
        echo Installation ist hier nicht moeglich.
        echo [HINWEIS] Nutze stattdessen eine gehostete PostgreSQL-Datenbank
        echo ^(z.B. Neon, Supabase, Railway, eigener 64-Bit-Server^) und trage deren
        echo Verbindungs-URL als DATABASE_URL in .env ein.
    ) else if "!HAS_WINGET!"=="1" (
        set /p INSTALL_PG="PostgreSQL jetzt ueber winget installieren? (j/N): "
        if /i "!INSTALL_PG!"=="j" (
            echo [INFO] Installiere PostgreSQL ueber winget ...
            rem Die ID "PostgreSQL.PostgreSQL" ist mehrdeutig (mehrere Versionen
            rem teilen sich diese ID) und wird von winget nicht gefunden. Es muss
            rem eine versionsspezifische ID verwendet werden.
            winget install --id PostgreSQL.PostgreSQL.18 -e --accept-source-agreements --accept-package-agreements
            if !ERRORLEVEL! neq 0 (
                echo [FEHLER] PostgreSQL-Installation ueber winget fehlgeschlagen.
                echo Bitte manuell installieren: https://www.postgresql.org/download/
                echo oder eine gehostete Datenbank ^(Neon, Supabase, Railway, ...^) nutzen
                echo und DATABASE_URL in .env entsprechend setzen.
                pause
                exit /b 1
            )
            where psql >nul 2>nul
            if !ERRORLEVEL! neq 0 (
                echo [WARNUNG] winget meldet Erfolg, aber "psql" ist noch nicht im PATH.
                echo Dieses Fenster schliessen und start.bat erneut ausfuehren, damit die
                echo PATH-Aenderung wirksam wird.
            ) else (
                echo [OK] PostgreSQL wurde installiert.
            )
            echo [HINWEIS] Bitte Datenbank/Nutzer gemaess DATABASE_URL in .env anlegen,
            echo bevor der Server gestartet wird.
        ) else (
            echo [HINWEIS] Ohne PostgreSQL funktioniert der Server nicht.
            echo Stelle sicher, dass eine erreichbare Datenbank ^(z.B. remote^) existiert.
        )
    ) else (
        echo [HINWEIS] winget nicht verfuegbar. Bitte PostgreSQL manuell installieren
        echo oder eine erreichbare Remote-Datenbank in .env eintragen:
        echo https://www.postgresql.org/download/
    )
) else (
    echo [OK] PostgreSQL-Client gefunden.
)
echo.

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

rem --- Go-Modulabhaengigkeiten sicherstellen ---
echo [INFO] Pruefe Go-Modulabhaengigkeiten ...
go mod download
if !ERRORLEVEL! neq 0 (
    echo [FEHLER] Go-Abhaengigkeiten konnten nicht geladen werden.
    pause
    exit /b 1
)
echo [OK] Abhaengigkeiten sind vorhanden.
echo.

set /p RUN_SEED="Datenbank-Migrationen/Demo-Daten jetzt seeden? (j/N): "
if /i "%RUN_SEED%"=="j" (
    echo [INFO] Fuehre Seed aus...
    go run ./cmd/seed
    if !ERRORLEVEL! neq 0 (
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
