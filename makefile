SHELL := /bin/bash

tidy:
	go mod tidy
	go mod vendor

VERSION := 1.0

all: chars-api

chars-api:
	docker build \
		-f infra/docker/dockerfile.chars-api \
		-t chars-api:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

deps-upgrade:
	go get -u -v ./...
	go mod tidy
	go mod vendor


# ------------------------------------------------------
# Running with k8s/kind

KIND_CLUSTER := chars-cluster

kind-up:
	kind create cluster \
		--image kindest/node:v1.24.1 \
		--name $(KIND_CLUSTER) \
		--config infra/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=chars-system

kind-load:
	cd infra/k8s/kind/chars-pod; kustomize edit set image chars-api-image=chars-api:$(VERSION)
	kind load docker-image chars-api:$(VERSION) --name $(KIND_CLUSTER)

kind-apply:
	kustomize build infra/k8s/kind/chars-pod | kubectl apply -f -

kind-update: all kind-load kind-restart

kind-update-apply: all kind-load kind-apply

kind-restart:
	kubectl rollout restart deployment chars-pod 

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status-chars:
	kubectl get pods -o wide --watch --namespace=chars-system


kind-logs:
	kubectl logs -l app=chars --all-containers=true -f --tail=100 | go run app/tooling/logfmt/main.go


kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=chars

kind-context-sales:
	kubectl config set-context --current --namespace=chars-system
