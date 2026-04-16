@echo off
set PG_CTL=..\.local\pgsql16\pgsql\bin\pg_ctl.exe
set PG_DATA=..\.local\pgsql16\data
set PG_LOG=..\.local\pgsql16\data\postgres.log

echo Memeriksa status PostgreSQL...
"%PG_CTL%" status -D "%PG_DATA%" >nul 2>&1
if %errorlevel% neq 0 (
    echo PostgreSQL belum berjalan. Mencoba menghidupkan...
    "%PG_CTL%" start -D "%PG_DATA%" -l "%PG_LOG%" -w
    if %errorlevel% equ 0 (
        echo PostgreSQL berhasil dihidupkan.
    ) else (
        echo Gagal menghidupkan PostgreSQL! Periksa file log di %PG_LOG%.
    )
) else (
    echo PostgreSQL sudah berjalan.
)

echo.
echo Memulai aplikasi backend...
go run main.go
