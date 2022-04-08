BINARY_NAME=omm

all: build

build:
	go build -o ./.out/${BINARY_NAME} .

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
