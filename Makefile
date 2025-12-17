swagger:
	swag init -o internal/common/infra/api/swagger

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -1

test-coverage-html:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Open coverage.html in your browser"