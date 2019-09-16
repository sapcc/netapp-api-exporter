app=netapp-api-exporter
IMAGE=hub.global.cloud.sap/monsoon/${app}

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git rev-parse HEAD | head -c 7)

ifeq ($(BRANCH),master)
	VERSION := $(HASH)
else
	VERSION := $(BRANCH)-$(HASH)
endif

# VERSION:=v$(shell date -u +%Y%m%d%H%M%S)

#netapp-api-exporter:

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
	DEV=1 ./bin/${app}_dev -l localhost -w 30

.PHONY: clean
clean:
	rm -f bin/${app}_linux_amd64
	rm -f bin/${app}_darwin_amd64
	rm -f bin/${app}_dev
