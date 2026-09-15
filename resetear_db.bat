@echo off
setlocal enabledelayedexpansion

echo ========================================================
echo        GIS-LACIS - Script de Reseteo de Base de Datos
echo ========================================================
echo.
echo ADVERTENCIA: Esto borrara toda la base de datos local y la creara
echo desde cero usando el archivo "Modelo Base de datos.sql".
echo.
pause

REM Buscar PostgreSQL en C:\Program Files\PostgreSQL
set "PG_PATH="
for /d %%D in ("C:\Program Files\PostgreSQL\*") do (
    if exist "%%D\bin\initdb.exe" (
        set "PG_PATH=%%D\bin"
    )
)

if exist "pg_data\" (
    "!PG_PATH!\pg_ctl.exe" -D "pg_data" status > nul 2>&1
    if not errorlevel 1 (
        echo [INFO] Deteniendo base de datos (por si estaba encendida)...
        "!PG_PATH!\pg_ctl.exe" -D "pg_data" stop -m fast
        if errorlevel 1 (
            echo [ERROR] No se pudo detener PostgreSQL - cerrala manualmente e intenta de nuevo.
            pause
            exit /b 1
        )
    )

    echo [INFO] Eliminando la base de datos antigua...
    rmdir /s /q "pg_data"
    if exist "pg_data\" (
        echo [ERROR] No se pudo eliminar la base de datos anterior ^(pg_data^). Puede haber quedado un proceso de Postgres usando esos archivos.
        pause
        exit /b 1
    )
    echo [INFO] Base de datos borrada con exito.
    echo.
)

echo [INFO] Arrancando iniciar.bat para crearla de cero...
call iniciar.bat
