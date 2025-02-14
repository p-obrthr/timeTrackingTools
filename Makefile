tidy:
	go mod tidy

build:
	go build -o bin/timeTrackingTools main.go

run: build
	./bin/timeTrackingTools

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/timeTrackingTools-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/timeTrackingTools-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/timeTrackingTools-freebsd-386 main.go

all: tidy run
