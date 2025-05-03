@echo off
echo Starting development environment...

:: Check if MongoDB is running
echo Checking MongoDB connection...
docker ps | findstr mongo > nul
if %errorlevel% neq 0 (
    echo MongoDB container not found. Please start MongoDB first.
    exit /b 1
)

:: Build React frontend
echo Building React frontend...
cd web
npm run build
cd ..

:: Start Go server
echo Starting Go server on port 6060...
echo Access the application at http://localhost:6060
go run cmd/main.go