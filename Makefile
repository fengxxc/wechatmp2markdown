BUILD_ENV := CGO_ENABLED=0
APP=wechatmp2markdown
VERSION=v1.1.0

# linux or mac 环境编译
# make [cmd]
build-linux: clean
	${BUILD_ENV} GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o build/${APP}-${VERSION}_linux_amd64 main.go
build-osx: clean
	${BUILD_ENV} GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o build/${APP}-${VERSION}_osx_amd64 main.go
build-win64: clean
	${BUILD_ENV} GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o build/${APP}-${VERSION}_win64.exe main.go
build-win32: clean
	${BUILD_ENV} GOOS=windows GOARCH=386 go build -ldflags "-s -w" -o build/${APP}-${VERSION}_win32.exe main.go


# windows环境编译 需gcc，推荐安装tdm64-gcc
# mingw32-make [cmd]
win-build-linux: clean
	go env -w ${BUILD_ENV}
	go env -w GOOS=linux
	go env -w GOARCH=amd64
	go build -ldflags "-s -w" -o build/${APP}-${VERSION}_linux_amd64 main.go
win-build-osx: clean
	go env -w ${BUILD_ENV}
	go env -w GOOS=darwin
	go env -w GOARCH=amd64
	go build -ldflags "-s -w" -o build/${APP}-${VERSION}_osx_amd64 main.go
win-build-win64: clean
	go env -w ${BUILD_ENV}
	go env -w GOOS=windows
	go env -w GOARCH=amd64
	go build -ldflags "-s -w" -o build/${APP}-${VERSION}_win64.exe main.go
win-build-win32: clean
	go env -w ${BUILD_ENV}
	go env -w GOOS=windows
	go env -w GOARCH=386
	go build -ldflags "-s -w" -o build/${APP}-${VERSION}_win32.exe main.go

run:
	go run main.go

clean:
	go clean