build:
	go build -o bin/acle cmd/main.go
	chmod +x bin/acle

run: build
	./bin/acle -if ./test_data/sample.ios -acl_id 103

test:
	go test -v

coverage:
	go test -coverprofile=c.out
	go tool cover -func=c.out
	rm c.out
