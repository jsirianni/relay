WORK_DIR := $(shell pwd)

VERSION := $(shell cat version)

build-all: build-frontend build-forwarder clean

build-frontend: clean-frontend
	mkdir -p docker/frontend/stage
	ls | grep -v docker | xargs -I{} cp -r {} docker/frontend/stage
	cd docker/frontend && docker build . -t firefoxx04/relay-frontend:latest

build-forwarder: clean-forwarder
	mkdir docker/forwarder/stage
	ls | grep -v docker | xargs -I{} cp -r {} docker/forwarder/stage
	cd docker/forwarder && docker build . -t firefoxx04/relay-forwarder:latest

quick-all: quick-frontend quick-forwarder

quick-frontend:
	cd cmd/frontend/ &&	go build

quick-forwarder:
	cd cmd/forwarder/ && go build

test-all:
	go test ./...

lint-all:
	golint ./...

fmt-all:
	go fmt ./...

clean: clean-frontend clean-forwarder

clean-frontend:
	rm -rf docker/frontend/stage

clean-forwarder:
	rm -rf docker/forwarder/stage

prune-docker:
	docker system prune --force

push: push-all

push-all:
	docker push firefoxx04/relay-frontend
	docker push firefoxx04/relay-forwarder

release: build-all
	docker tag firefoxx04/relay-frontend:latest firefoxx04/relay-frontend:${VERSION}
	docker tag firefoxx04/relay-forwarder:latest firefoxx04/relay-forwarder:${VERSION}
	docker push firefoxx04/relay-frontend:${VERSION}
	docker push firefoxx04/relay-forwarder:${VERSION}
