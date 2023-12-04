BINARY_NAME=warp

all: build

build:
	go build -o bin/$(BINARY_NAME) main.go

clean:
	rm -rf bin/