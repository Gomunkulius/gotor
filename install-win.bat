@echo off

curl https://github.com/Gomunkulius/gotor/releases/latest/download/gotor-win-amd64.exe -O

echo "Renaming"
Rename ".\gotor-win-amd64.exe""gotor.exe"
echo "Creating directory"
mkdir C:/Program Files/gotor
Move .\gotor.exe C:/Program Files/gotor