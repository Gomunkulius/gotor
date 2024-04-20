@echo off

curl https://github.com/Gomunkulius/gotor/releases/latest/download/gotor-win-amd64.exe -OutFile gotor-win-amd64.exe

echo "Renaming"
rename ".\gotor-win-amd64.exe" "gotor.exe"
echo "Creating directory"
mkdir "C:/Program Files/gotor"
move ".\gotor.exe" "C:/Program Files/gotor"
pause