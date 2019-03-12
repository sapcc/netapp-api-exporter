ARCH=amd64
OS=linux
IMAGE=mx3d/netapp-api-exporter
VERSION=v0.1

build:
	GOOS=$(OS) GOARCH=$(ARCH) go build
	docker build -t $(IMAGE):$(VERSION) . 

push:
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

test:
	GOOS=darwin GOARCH=$(ARCH) go build
