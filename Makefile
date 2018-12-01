BINARY=ratelimit

.PHONY=build
build:
	go build -o ${BINARY} main.go

${BINARY}: build
