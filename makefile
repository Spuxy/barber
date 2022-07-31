VERSION := 1
COMMIT?=$(shell git rev-parse --short HEAD)
DATE=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
APP_NAME="barber"

#==============================================================
#GOLANG Stuffs
run:
	go run main.go -ldflags "-X main.buildVersion=$(VERSION) -X main.buildCommit=$(COMMIT) -X main.buildDate=($DATE)"

build:
	go build -ldflags "-X main.buildVersion=$(VERSION) -X main.buildCommit=$(COMMIT) -X main.buildDate=($DATE)" -o "$(APP_NAME)-amd64"

tidy:
	go mod tidy && go mod vendor

clean:
	go clean
	rm -rf "$(APP_NAME)-amd64"

#==============================================================
#Docker Stuffs
docker-build:
	docker build -f ./.docker/Dockerfile \
	-t barber-api-amd64:$(VERSION) \
	--build-arg BUILD_VERSION=$(VERSION) \
	--build-arg BUILD_COMMIT=$(COMMIT) \
	--build-arg BUILD_DATE=$(DATE) .

#==============================================================
# K8s Stuffs
KIND_NAMESPACE := barber-system
KIND_CLUSTER := barber-system-cluster
kind-up:
	kind create cluster \
		--image kindest/node:v1.24.0@sha256:0866296e693efe1fed79d5e6c7af8df71fc73ae45e3679af05342239cdc5bc8e \
		--name $(KIND_CLUSTER) \
		--config .k8s/kind/kind-config.yml
	kubectl config set-context --current --namespace=$(KIND_NAMESPACE)

kind-load:
	kind load docker-image barber-api-amd64:$(VERSION) --name $(KIND_CLUSTER)
	cd zarf/k8s/kind/sales-pod; kustomize edit set image sales-api-image=sales-api-amd64:$(VERSION)