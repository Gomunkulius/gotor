Invoke-WebRequest https://github.com/Gomunkulius/gotor/releases/latest/download/gotor-win-amd64.exe -OutFile gotor-win-amd64.exe

echo "Creating directory"
mkdir C:/Program Files/gotor
Move-Item -Path .\gotor-win-amd64.exe -Destination C:/Program Files/gotor
