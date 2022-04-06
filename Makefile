.PHONY: install
install:
	@go mod tidy
	@go install \
	  github.com/abice/go-enum \
	  github.com/golang/mock/mockgen

.PHONY: gen
gen: install
	@go generate ./...

.PHONY: test
test: gen
	@go test $$(go list ./... | grep -v /mock) -coverprofile=coverage.txt -covermode=atomic -timeout 60s

.PHONY: coverage
coverage:
	@make test | tr -d '\n' | sed -e 's/.*coverage:\(.*\)of statements.*/\1/' | tr -d '\n '

.PHONY: lint
lint:
	@golangci-lint --timeout 10m0s run
