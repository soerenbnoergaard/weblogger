.PHONY: all build install

TARGET = github.com/soerenbnoergaard/weblogger

all: build

build:
	go build $(TARGET)

install:
	go install $(TARGET)
