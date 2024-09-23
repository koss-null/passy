.PHONY: all build clean

all: clean build

build:
	go build -o build/passy ./

test: build
	./build/passy -h
	./build/passy -c -readable
	./build/passy -c -safe
	./build/passy -c -insane

clean:
	rm -rf build/passy
