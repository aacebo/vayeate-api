clean:
	rm -rf ./bin

build:
	go build -o bin/main main.go

clean_build: clean build

run:
	go run ./...

fmt:
	gofmt -w .

test:
	go clean -testcache
	go test ./... -cover

compile:
	GOOS=freebsd GOARCH=386 go build -o bin/main-freebsd-386 main.go
	GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
	GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
