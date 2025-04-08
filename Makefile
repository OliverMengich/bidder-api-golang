delete:
	rm ./bidder-api-golang.exe
build:
	go build
start:
	./bidder-api-golang.exe
restart: delete build start