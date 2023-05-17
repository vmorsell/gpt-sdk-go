.PHONY: test
test:
	go test -short -v  ./...

.PHONY: testall
testall:
	go test -v ./...
