.PHONY: build xbuild build-windows build-linux

build:
	GOOS=darwin GOARCH=amd64 go build -o bin/imageinfo

clean:
	rm -fr bin

xbuild: clean build build-windows build-linux
	zip bin/imageinfo.darwin-amd64.zip bin/imageinfo

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/imageinfo
	zip bin/imageinfo.linux-amd64.zip bin/linux/imageinfo

build-windows:
	GOOS=windows GOARCH=amd64 go build -o bin/windows/imageinfo
	zip bin/imageinfo.windows-amd64.zip bin/windows/imageinfo
