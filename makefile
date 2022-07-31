VERSION := 1
COMMIT?=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
APP_NAME="barber"

run:
	go run main.go -ldflags "-X main.buildVersion=$(VERSION) -X main.buildCommit=$(COMMIT) -X main.buildDate=($DATE)"

build:
	go build -ldflags "-X main.buildVersion=$(VERSION) -X main.buildCommit=$(COMMIT) -X main.buildDate=($DATE)" -o "$(APP_NAME)-amd64"

tidy:
	go mod tidy && go mod vendor

clean:
	go clean
	rm -rf "$(APP_NAME)-amd64"

docker-build:
	docker build -f ./.docker/Dockerfile \
	-t barber-api-amd64:$(VERSION) \
	--build-arg BUILD_VERSION=$(VERSION) \
	--build-arg BUILD_COMMIT=$(COMMIT) \
	--build-arg BUILD_DATE=$(DATE) .