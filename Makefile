.PHONY: test
test:
	go test .

.PHONY: test-coverage
test-coverage:
	mkdir -p tmp
	go test -coverprofile=./tmp/coverage.out .
	go tool cover -html=tmp/coverage.out
	open tmp/coverage.html

.PHONY: clean
clean:
	rm tmp/*
	go clean -testcache
	go clean -modcache
	go clean -cache
