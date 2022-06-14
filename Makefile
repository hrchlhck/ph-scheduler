all: build run

build:
	go build -o scheduler .

run:
	./scheduler teste
