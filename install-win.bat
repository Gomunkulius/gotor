@echo off
echo Downloading latest release...
Powershell.exe -command "Invoke-WebRequest https://github.com/Gomunkulius/gotor/releases/latest/download/gotor-win-amd64.exe -OutFile gotor-win-amd64.exe"

echo Renaming..
rename ".\gotor-win-amd64.exe" "gotor.exe"
echo Creating directory...
mkdir "C:/Program Files/gotor"
move ".\gotor.exe" "C:\Program Files\gotor"
where 7z.exe 2>nul || setx PATH "%PATH%;C:\Program Files\gotor" /M
pause
