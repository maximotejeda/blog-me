export DEBUG=true
export SERVERADDR=
export SERVERPORT=8080
export FILESDIR=./posts
export ASSETSDIR=./assets
export DBDIR=./db
export DBFILE=blog.db

PHONY: all, clean, build, container, push, run, assets

build:
	go build -o bin/blog-me cmd/main.go
	
run: build
	./bin/blog-me	
image: 
	docker build -t localhost:32000/blog-me .
deploy: 
	ssh w1 mkdir ~/golang/blog-me 2>/dev/null | ls 1> /dev/null
	rsync -r ./ w1:~/golang/blog-me
	ssh -C w1 "cd ~/golang/blog-me/ && make image && docker push localhost:32000/blog-me"

copystatic:
	rsync -rv ./assets/* w1:/srv/ext/nfsShare/kubernetes/volumes/pvc-81221464-1ff2-4340-99da-976b59acedba
copypost:
	rsync -rv ./posts/*.html w1:/srv/ext/nfsShare/kubernetes/volumes/pvc-22e3bc14-8660-4703-b2d3-055b37dcf74d
test:
	go test ./...
clean:
	rm -rf db bin 
	@echo "Cleaning project"
