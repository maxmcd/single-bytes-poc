

build: Dockerfile
	docker build -t single-bytes-poc .

run: build
	docker run -it -v $$(pwd):/opt single-bytes-poc bash
