export DEBUG=true
export FILESDIR=./files
export ASSETSDIR=./posts/assets

PHONY: all, clean, build, container, push, run

build:
	go build -o bin/blog-me cmd/main.go
run: build
	./bin/blog-me	

container: build
	docker build -t localhost:32000/blog-me .


clean:
	@echo "Cleaning project"
