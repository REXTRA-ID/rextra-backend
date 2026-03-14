@echo off
TITLE REXTRA API TESTER
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "%~dp0test_transaction.ps1"
pause
