app=netapp-api-exporter
IMAGE=keppel.eu-de-1.cloud.sap/ccloud/${app}

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
HASH := $(shell git log -1 --format=%h)
TAG :=$(shell date -u +%Y%m%d%H%M%S)-$(BRANCH)-$(HASH)

GOFILES := $(wildcard *.go) $(wildcard pkg/**/*.go)

DEV_ARGS := -l localhost
ifneq ($(strip $(CONFIG_FILE)),)
DEV_ARGS += -c $(CONFIG_FILE)
endif

.PHONY: build
build: bin/${app}_linux_amd64 bin/${app}_darwin_amd64

bin/${app}_linux_amd64: $(GOFILES)
	GOOS=linux GOARCH=amd64 go build -o $@

bin/${app}_darwin_amd64: $(GOFILES)
	GOOS=darwin GOARCH=amd64 go build -o $@

bin/${app}_darwin_arm64: $(GOFILES)
	GOOS=darwin GOARCH=arm64 go build -o $@

.PHONY: docker
docker: bin/${app}_linux_amd64
	@echo "[INFO] build docker image"
	docker build --platform=linux/amd64 --progress=plain --no-cache -t $(IMAGE):$(TAG) .
	@echo "[INFO] push docker image"
	docker tag $(IMAGE):$(TAG) $(IMAGE):latest
	docker push $(IMAGE):$(TAG)
	docker push $(IMAGE):latest

.PHONY: dev
dev: *.go
	DEV=1 go run $^ $(DEV_ARGS)

.PHONY: clean
clean:
	rm -f bin/*
