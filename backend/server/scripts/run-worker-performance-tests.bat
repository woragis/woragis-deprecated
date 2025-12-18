@echo off
REM Script to run performance tests for all workers on Windows
REM Usage: run-worker-performance-tests.bat [worker-name]

setlocal enabledelayedexpansion

set SCRIPT_DIR=%~dp0
set BACKEND_DIR=%SCRIPT_DIR%..\..
set SERVER_DIR=%BACKEND_DIR%\server

echo Starting Worker Performance Tests
echo ==========================================

REM Check if Docker Compose is available
where docker-compose >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo Error: docker-compose not found
    exit /b 1
)

REM Start test dependencies
echo.
echo Starting test dependencies...
cd /d "%SERVER_DIR%"
docker-compose -f docker-compose.test.yml up -d

REM Wait for services to be ready
echo Waiting for services to be ready...
timeout /t 5 /nobreak >nul

REM Function to run performance tests for a worker
if "%1"=="" (
    REM Run all worker performance tests
    call :run_worker_tests email-worker
    call :run_worker_tests translation-worker
    call :run_worker_tests whatsapp-worker
) else (
    REM Run specific worker tests
    call :run_worker_tests %1
)

REM Cleanup
echo.
echo Cleaning up...
cd /d "%SERVER_DIR%"
docker-compose -f docker-compose.test.yml down

echo.
echo Performance tests completed!
exit /b 0

:run_worker_tests
set WORKER_NAME=%1
set WORKER_DIR=%BACKEND_DIR%\%WORKER_NAME%

if not exist "%WORKER_DIR%" (
    echo Error: Worker directory not found: %WORKER_DIR%
    exit /b 1
)

echo.
echo Running performance tests for %WORKER_NAME%...
echo ----------------------------------------

cd /d "%WORKER_DIR%"

if exist "internal\integration" (
    REM Run performance tests (skip short tests)
    go test -tags=integration -run "Test.*Performance|Test.*Load|Test.*Latency|Test.*Rate|Benchmark" ./internal/integration/... -v
    
    if %ERRORLEVEL% EQU 0 (
        echo [OK] %WORKER_NAME% performance tests passed
    ) else (
        echo [FAIL] %WORKER_NAME% performance tests failed
        exit /b 1
    )
) else (
    echo No integration tests found for %WORKER_NAME%
)

exit /b 0
