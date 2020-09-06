all: lint test

install:
	@echo "Installing dev dependencies..."
	@(curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin" v1.27.0)
	@GO111MODULE=off go get -u github.com/joho/godotenv/cmd/godotenv

lint:
	@golangci-lint run

todo:
	@golangci-lint run --no-config --disable-all --enable godox

test:
	@go test ./...

integration_test:
	@godotenv -f .env go test -v ./... -tags=integration