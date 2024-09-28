.PHONY: all build clean

all: clean build

build:
	go build -o build/passy ./

clean:
	rm -rf build/passy
