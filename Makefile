

build: Dockerfile
	docker build -t single-bytes-poc .

run: build
	docker run -v $$(pwd):/opt single-bytes-poc
