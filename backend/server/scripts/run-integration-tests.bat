@echo off
REM Script to run integration tests on Windows
REM Usage: scripts\run-integration-tests.bat

echo Starting test dependencies...
docker-compose -f docker-compose.test.yml up -d

echo Waiting for services to be healthy...
timeout /t 10 /nobreak

echo Checking service health...
docker-compose -f docker-compose.test.yml ps

echo Running integration tests...
go test ./app/internal/integration/... -tags=integration -v

echo Cleaning up...
docker-compose -f docker-compose.test.yml down

echo Integration tests completed!
