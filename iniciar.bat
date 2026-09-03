@echo off
setlocal enabledelayedexpansion

echo ========================================================
echo        GIS-LACIS - Script de Inicio Automtico
echo ========================================================
echo.

echo [INFO] Cerrando sesiones de PostgreSQL que hayan quedado abiertas...
taskkill /F /IM postgres.exe /T > nul 2>&1

echo [INFO] Cerrando el servidor web que haya quedado abierto (puerto 8080)...
for /f "tokens=5" %%P in ('netstat -ano ^| findstr ":8080" ^| findstr "LISTENING"') do (
    taskkill /F /PID %%P /T > nul 2>&1
)

timeout /t 1 /nobreak > nul

REM Buscar PostgreSQL en C:\Program Files\PostgreSQL
set "PG_PATH="
for /d %%D in ("C:\Program Files\PostgreSQL\*") do (
    if exist "%%D\bin\initdb.exe" (
        set "PG_PATH=%%D\bin"
    )
)

if "!PG_PATH!"=="" (
    echo [ERROR] No se encontro PostgreSQL instalado en C:\Program Files\PostgreSQL.
    echo Por favor, instala PostgreSQL 18 o superior.
    pause
    exit /b 1
)

echo [INFO] PostgreSQL encontrado en: !PG_PATH!

REM Verificar si el cluster ya existe
if not exist "pg_data\" (
    echo [INFO] Inicializando nuevo cluster de base de datos local...
    "!PG_PATH!\initdb.exe" -D "pg_data" -U postgres --auth=trust > nul
    
    REM Cambiar puerto a 5433 en postgresql.conf
    (
        echo port = 5433
    ) >> "pg_data\postgresql.conf"

    echo [INFO] Iniciando servidor de base de datos...
    start /B "" "!PG_PATH!\pg_ctl.exe" -D "pg_data" -l "pg_data\pg.log" start > nul
    
    echo Esperando a que el servidor inicie...
    timeout /t 3 /nobreak > nul

    echo [INFO] Creando base de datos 'lacis' y tablas...
    "!PG_PATH!\psql.exe" -p 5433 -U postgres -d postgres -c "ALTER USER postgres WITH PASSWORD 'isma_mesa22';" > nul
    "!PG_PATH!\psql.exe" -p 5433 -U postgres -d postgres -c "CREATE DATABASE lacis;" > nul
    "!PG_PATH!\psql.exe" -p 5433 -U postgres -d lacis -f "Modelo Base de datos\Modelo Base de datos.sql" > nul
    
    echo [INFO] Base de datos creada exitosamente.
) else (
    echo [INFO] El cluster de base de datos ya existe. Iniciando servidor...
    start /B "" "!PG_PATH!\pg_ctl.exe" -D "pg_data" -l "pg_data\pg.log" start > nul
    timeout /t 2 /nobreak > nul
)

echo.
echo ========================================================
echo [INFO] Iniciando el Servidor Web (Go)...
echo La pagina estara disponible en: http://localhost:8080
echo ========================================================
echo.
start "" "cmd.exe" /k "go run cmd/api/main.go"

echo.
echo Presiona cualquier tecla para detener la base de datos...
pause > nul

echo [INFO] Deteniendo base de datos...
"!PG_PATH!\pg_ctl.exe" -D "pg_data" stop > nul
taskkill /F /IM postgres.exe /T > nul 2>&1

