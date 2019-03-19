ARCH=amd64
OS=linux
IMAGE=mx3d/netapp-api-exporter
VERSION=v0.1

netapp-api-exporter: *.go
	@echo "[INFO] build go excutable for $(ARCH)"
	GOOS=$(OS) GOARCH=$(ARCH) go build

.PHONY: build
build: netapp-api-exporter
	@echo "[INFO] build docker image"
	docker build -t $(IMAGE):$(VERSION) . 

.PHONY: push
push: build
	@echo "[INFO] push docker image"
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

.PHONY: dev
dev: 
	go build
	DEV=1 ./netapp-api-exporter -l localhost -w 30