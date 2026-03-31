@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

set "DB_HOST=%DB_HOST%:localhost"
set "DB_PORT=%DB_PORT%:3306"
set "DB_USER=%DB_USER%:root"
set "DB_PASSWORD=%DB_PASSWORD%:123456"
set "DB_NAME=%DB_NAME%:bycigar"

set "SCRIPT_DIR=%~dp0"
cd /d "%SCRIPT_DIR%"

echo ==========================================
echo BYCIGAR 数据库初始化
echo ==========================================
echo.

set /p ENV="输入环境 (dev/test，默认 dev): "
if "%ENV%"=="" set "ENV=dev"
if "%ENV%"=="test" set "DB_NAME=bycigar_test"

echo.
echo 数据库配置:
echo   地址: %DB_HOST%
echo   端口: %DB_PORT%
echo   用户: %DB_USER%
echo   数据库: %DB_NAME%
echo.

echo [1/5] 检查MySQL连接...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% -e "SELECT 1" >nul 2>&1
if errorlevel 1 (
    echo [错误] 无法连接到MySQL，请检查配置
    pause
    exit /b 1
)
echo [OK] MySQL连接正常

echo.
echo [2/5] 创建数据库...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% -e "CREATE DATABASE IF NOT EXISTS `%DB_NAME%` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" >nul 2>&1
echo [OK] 数据库创建完成

echo.
echo [3/5] 导入表结构 (001_init_schema.sql)...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% %DB_NAME% < "001_init_schema.sql" >nul 2>&1
if errorlevel 1 (
    echo [错误] 表结构导入失败
    pause
    exit /b 1
)
echo [OK] 表结构导入完成

echo.
echo [4/5] 导入种子数据...
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% %DB_NAME% < "002_seed_base.sql" >nul 2>&1
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% %DB_NAME% < "003_seed_demo.sql" >nul 2>&1
mysql -h%DB_HOST% -P%DB_PORT% -u%DB_USER% -p%DB_PASSWORD% %DB_NAME% < "004_seed_settings.sql" >nul 2>&1
echo [OK] 种子数据导入完成

echo.
echo ==========================================
echo 数据库初始化完成！
echo ==========================================
echo.
echo 管理员账号:
echo   Email: admin@bycigar.com
echo   Password: 123456
echo.
echo 客服账号:
echo   Email: service1@bycigar.com
echo   Password: 123456
echo.

pause
