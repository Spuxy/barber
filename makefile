VERSION := 1
COMMIT?=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

run:
	go build -ldflags "-X main.buildVersion=build"

tidy:
	go mod tidy && go mod vendor

docker-build:
	docker build -f ./.docker/Dockerfile \
	-t barber-api-amd64:$(VERSION) \
	--build-arg BUILD_VERSION=$(VERSION) \
	--build-arg BUILD_COMMIT=$(COMMIT) \
	--build-arg BUILD_DATE=$(DATE) .