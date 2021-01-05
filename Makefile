DOCKER_IMAGE_VER=0.1
DOCKER_IMAGE=webhook-demo:${DOCKER_IMAGE_VER}

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' ./cmd/initc
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' ./cmd/webhook

docker-build:
	docker build -t ${DOCKER_IMAGE} .

docker-push:
	docker tag ${DOCKER_IMAGE} docker.io/dmitsh/${DOCKER_IMAGE} && docker push docker.io/dmitsh/${DOCKER_IMAGE}

.PHONY: build docker-build docker-push
