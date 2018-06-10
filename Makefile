all: chimaera

godep:
	@go get -u github.com/golang/dep/...

deps: godep
	@dep ensure

chimaera: deps
	@go build -o chimaera .

clean:
	@rm -rf chimaera

docker:
	@docker build -t chimaera:latest .
