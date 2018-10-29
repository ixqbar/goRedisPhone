SHELL:=/bin/bash
TARGET=goRedisPhone

all: win linux mac

win: 
	GOOS=windows GOARCH=amd64 go build -o ./bin/${TARGET}.exe ./src
	GOOS=windows GOARCH=386 go build -o ./bin/${TARGET}-x86.exe ./src

linux: 
	GOOS=linux GOARCH=amd64 go build -o ./bin/${TARGET}_${@} ./src

mac: 
	GOOS=darwin GOARCH=amd64 go build -o ./bin/${TARGET}_${@} ./src
	
clean:
	rm -rf ./bin/${TARGET}*