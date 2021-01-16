DOCKER_REPO=docker.io/dmitsh
DOCKER_IMAGE_VER=0.1
DOCKER_IMAGE=webhook-demo:${DOCKER_IMAGE_VER}

build:
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' ./cmd/initc
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' ./cmd/webhook

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' ./cmd/initc
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' ./cmd/webhook

docker-build:
	docker build -t ${DOCKER_IMAGE} .

docker-build-slim:
	docker build -t ${DOCKER_IMAGE} -f ./Dockerfile-slim .

docker-push:
	docker tag ${DOCKER_IMAGE} ${DOCKER_REPO}/${DOCKER_IMAGE} && docker push ${DOCKER_REPO}/${DOCKER_IMAGE}

.PHONY: build build-linux docker-build docker-build-slim docker-push
