run:
	go run cmd/main/main.go

install:
	go build -o gotor cmd/main/main.go
	mkdir /usr/bin/gotor
	mv ./gotor /usr/bin/gotor