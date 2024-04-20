command -v go >/dev/null 2>&1 || { echo >&2 "Golang is not installed. Please install it"; exit 1; }

echo "Installing gotor"
echo "Building binary..."
go build -o gotor cmd/main/main.go
echo "Installing binary..."
mv ./gotor /usr/local/bin
echo "Installation complete"
