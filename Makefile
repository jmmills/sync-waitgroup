dep:
	go mod tidy

test: 
	go test -race -v -coverprofile cover.out ./...

cover.html: test
	go tool cover -html cover.out -o cover.html

coverage: cover.html
	open cover.html

lint:
	gofmt -w . 
	golint 
	go vet

clean:
	rm -v cover.out cover.html
