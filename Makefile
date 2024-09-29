.PHONY: all build clean

all: clean build

build:
	go build -o build/passy ./

gccgobuild: main.go
	go build -gcflags="-B -C" -ldflags="-s -w" -o build/passy ./

clean:
	rm -rf build/passy
