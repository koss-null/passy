.PHONY: all build clean

all: clean build

build:
	go build -o build/passy ./

run: build
	./build/passy -i

build-gccgo: main.go
	go build -gcflags="-B -C" -ldflags="-s -w" -o build/passy ./

release-build-release: build-other-platforms
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-s -w -extldflags=-static" \
		-trimpath \
		-o build/passy ./
	upx --best --lzma build/passy

build-other-platforms:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
		-ldflags="-s -w -extldflags=-static" \
		-trimpath \
		-o build/passy.exe
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build \
		-ldflags="-s -w -extldflags=-static" \
		-trimpath \
		-o build/passy_osx

clean:
	rm -rf build/passy
