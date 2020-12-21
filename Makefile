SRCS = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

BUILD_FLAGS = -ldflags ""

.DEFAULT_GOAL := build/covid-alert

.PHONY: clean build-arm covid-alert-builder

build:
	@mkdir build

build/covid-alert: build $(SRCS)
	@go build ${BUILD_FLAGS} -o build/covid-alert .

build/covid-alert-arm: build $(SRCS)
	@CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 go build -o build/covid-alert-arm .

build-arm:
	@docker run -it --rm -v `pwd`:/workspace covid-alert-builder

covid-alert-builder:
	@cd builder && docker build -t covid-alert-builder .

clean:
	@rm -rf build