run:
	go run cmd/main/main.go

build:
	go build -o gotor cmd/main.go

install:
	go build -o gotor cmd/main/main.go
	mv ./gotor /usr/local/bin