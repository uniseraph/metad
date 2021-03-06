# Copyright (c) 2015-2017, ANT-FINANCE CORPORATION. All rights reserved.

SHELL = /bin/bash

GOLANG = golang:1.8.3

PROJECT = github.com/zanecloud/metad
TARGET  = metad
VERSION = $(shell cat VERSION)
GITCOMMIT = $(shell git log -1 --pretty=format:%h)
BUILD_TIME = $(shell date --rfc-3339 ns 2>/dev/null | sed -e 's/ /T/')

IMAGE_NAME = registry.cn-hangzhou.aliyuncs.com/zanecloud/metad

build:
	docker run -v $(shell pwd):/go/src/${PROJECT} -w /go/src/${PROJECT} --rm ${GOLANG} make local

binary: build

local:
	rm -rf bundles/${VERSION}
	mkdir -p bundles/${VERSION}/binary
	CGO_ENABLED=0 go build -v -ldflags "-X main.Version=${VERSION} -X main.GitCommit=${GITCOMMIT} -X main.BuildTime=${BUILD_TIME}" -o bundles/${VERSION}/binary/metad

image:build
	docker build -t ${IMAGE_NAME} .
	docker tag ${IMAGE_NAME} ${IMAGE_NAME}:${VERSION}-${GITCOMMIT}
	docker tag ${IMAGE_NAME} ${IMAGE_NAME}:${VERSION}

run:local
	chmod +x bundles/${VERSION}/binary/metad
	./bundles/${VERSION}/binary/metad  --log-level debug start --consul-addr 127.0.0.1:8500 --addr 127.0.0.1:6400

publish:
	docker tag ${IMAGE_NAME}:${VERSION}-${GITCOMMIT} ${IMAGE_NAME}:${VERSION}
	docker tag ${IMAGE_NAME}:${VERSION}-${GITCOMMIT} ${IMAGE_NAME}
	docker push ${IMAGE_NAME}:${VERSION}-${GITCOMMIT}
	docker push ${IMAGE_NAME}:${VERSION}
	docker push ${IMAGE_NAME}

.PHONY: build binary local image release
