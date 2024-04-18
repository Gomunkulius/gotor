run:
	go run cmd/main/main.go

build:
	go build -o gotor cmd/main.go

test:
	go test ./...

build_all_platforms:
	mkdir build
	cd build
	go env -w GOOS='windows'
	go env -w GOARCH="amd64"
	go env -w CGO_ENABLED='0'
	sudo go build -o build/gotor-win-amd64.exe cmd/main/main.go
	echo "Builded for win amd64"
	go env -w GOOS='linux'
	go env -w GOARCH="arm"
	GOOS=linux GOARCH=arm CGO_ENABLED=0 sudo go build -o build/gotor-linux-arm cmd/main/main.go
	echo "Builded for linux arm"
	go env -w GOOS='linux'
	go env -w GOARCH="amd64"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 sudo go build -o build/gotor-linux-amd64 cmd/main/main.go
	echo "Builded for linux amd64"
	go env -w GOOS='darwin'
	go env -w GOARCH="arm64"
	GOOS=darwin GOARCH=arm CGO_ENABLED=0 sudo go build -o build/gotor-darwin-arm cmd/main/main.go
	echo "Builded for darwin arm"
	go env -w GOOS='darwin'
	go env -w GOARCH="amd64"
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 sudo go build -o build/gotor-darwin-amd64 cmd/main/main.go
	echo "Builded for darwin amd64"
	go env -w GOOS='darwin'
	go env -w GOARCH="arm64"
	go env -w CGO_ENABLED='1'

install:
	go build -o gotor cmd/main/main.go
	mv ./gotor /usr/local/bin