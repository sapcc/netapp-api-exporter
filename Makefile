app=netapp-api-exporter
IMAGE=keppel.eu-de-1.cloud.sap/ccloud/${app}

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse HEAD | head -c 7)
VERSION :=$(BRANCH)-$(HASH)-$(shell date -u +%Y%m%d%H%M%S)

.PHONY: build
build: bin/${app}_linux_amd64 bin/${app}_darwin_amd64

bin/${app}_linux_amd64: *.go
	GOOS=linux GOARCH=amd64 go build -o $@

bin/${app}_darwin_amd64: *.go
	GOOS=darwin GOARCH=amd64 go build -o $@

.PHONY: docker
docker: bin/${app}_linux_amd64
	@echo "[INFO] build docker image"
	docker build -t $(IMAGE):$(VERSION) .
	@echo "[INFO] push docker image"
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
	docker push $(IMAGE):$(VERSION)
	docker push $(IMAGE):latest

.PHONY: dev
dev: 
	rm -f bin/${app}_dev
	go build -o bin/${app}_dev
	DEV=1 ./bin/${app}_dev -c config/netapp_filers.yaml -l localhost

.PHONY: clean
clean:
	rm -f bin/*
