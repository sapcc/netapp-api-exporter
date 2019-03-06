ARCH=amd64
OS=linux
IMAGE=mx3d/netapp-api-exporter
VERSION=v0.1

build:
	GOOS=$(OS) GOARCH=$(ARCH) go build
	docker build --build-arg OS_USERNAME=${OS_USERNAME} --build-arg OS_PASSWORD=${OS_PASSWORD} -t $(IMAGE):$(VERSION) . 

push:
	docker push $(IMAGE):$(VERSION)

test:
	GOOS=darwin GOARCH=$(ARCH) go build
