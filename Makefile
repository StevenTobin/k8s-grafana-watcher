WORKING_DIR = $(shell pwd)
IMAGE_NAME = kube-grafana-watcher
VERSION = 0.1

default: dep build docker_image

# Ensure vendor has the dependencies
dep:
	docker run \
            --rm -i \
            -v $(WORKING_DIR):/go/src/github.com/PierreVincent/grafana-watcher \
            -w /go/src/github.com/PierreVincent/grafana-watcher \
            supinf/go-dep ensure

# Build executable
build:
	docker run \
        --rm -i \
        -e GOARCH=amd64 \
        -e CGO_ENABLED=0 \
        -v $(WORKING_DIR):/home/stobin/repos/go-workspace/src/src/github.com/k8s-grafana-watcher \
        -w /home/stobin/repos/go-workspace/src/src/github.com/k8s-grafana-watcher \
        golang:1.8 go build -ldflags="-s -w" -v -o bin/main *.go

# Build docker image
docker_image:
	docker build -t $(IMAGE_NAME):$(VERSION) .