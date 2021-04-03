GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=peppermint-server

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME)
test:
	$(GOTEST) ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
run: build
	MONGODB_USERNAME=root MONGODB_PASSWORD=example MONGODB_ENDPOINT=127.0.0.1:27017 WEBSERVER_PORT=1337 ./$(BINARY_NAME)
testcoverage:
	$(GOTEST) -coverprofile coverage.out ./... && go tool cover -html=coverage.out && rm coverage.out
