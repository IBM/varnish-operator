# Image URL to use in all building/pushing image targets
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION ?= local
REPO ?= cinple
PUBLISH_IMG ?= ${REPO}/varnish-operator:${VERSION}
IMG ?= ${PUBLISH_IMG}-dev
VARNISH_PUBLISH_IMG ?= varnish:${VERSION}
VARNISH_IMG ?= ${VARNISH_PUBLISH_IMG}-dev
VARNISH_CONTROLLER_PUBLISH_IMG ?= varnish-controller:${VERSION}
VARNISH_CONTROLLER_IMG ?= ${VARNISH_CONTROLLER_PUBLISH_IMG}-dev
VARNISH_METRICS_PUBLISH_IMG ?= varnish-metrics-exporter:${VERSION}
VARNISH_METRICS_IMG ?= ${VARNISH_METRICS_PUBLISH_IMG}-dev
NAMESPACE ?= "default"
CRD_OPTIONS ?= "crd:crdVersions=v1"
PLATFORM ?= "linux/amd64"

# all: test varnish-operator
all: test varnish-operator varnish-controller

# Run tests
test: generate fmt vet manifests
	go test github.com/cin/varnish-operator/pkg/... github.com/cin/varnish-operator/cmd/... github.com/cin/varnish-operator/api/... -coverprofile cover.out

# Run lint tools
lint:
	golangci-lint run

# Build varnish-operator binary
varnish-operator: generate fmt vet
	go build -o ${ROOT_DIR}bin/varnish-operator github.com/cin/varnish-operator/cmd/varnish-operator

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	NAMESPACE=${NAMESPACE} LOGLEVEL=debug LOGFORMAT=console CONTAINER_IMAGE=${REPO}/${VARNISH_IMG} LEADERELECTION_ENABLED=false WEBHOOKS_ENABLED=false go run ${ROOT_DIR}cmd/varnish-operator/main.go

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
install: manifests
	kustomize build ${ROOT_DIR}config/crd | kubectl apply -f -

uninstall:
	kustomize build ${ROOT_DIR}config/crd | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	# CRD apiextensions.k8s.io/v1
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=varnish-operator paths="./..." output:crd:artifacts:config=config/crd/bases
	kustomize build ${ROOT_DIR}config/crd > $(ROOT_DIR)varnish-operator/crds/varnishcluster.yaml

	# ClusterRole
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=varnish-operator paths="./..." output:crd:none output:rbac:stdout > $(ROOT_DIR)varnish-operator/templates/clusterrole.yaml

# Run goimports against code
fmt:
	cd ${ROOT_DIR} && goimports -w ./pkg ./cmd ./api

# Run go vet against code
vet:
	cd ${ROOT_DIR} && go vet ./pkg/... ./cmd/... ./api/...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile=./hack/boilerplate.go.txt paths="./..."

helm-upgrade: manifests
	helm upgrade --install --namespace ${NAMESPACE} --force varnish-operator --set operator.controllerImage.tag=${VERSION} --set namespace=${NAMESPACE} ${ROOT_DIR}varnish-operator

# Build the docker image
# docker-build: test
docker-build: test
	docker build --platform ${PLATFORM} ${ROOT_DIR} -t ${IMG} -f Dockerfile

# Tag and push the docker image
docker-tag-push:
ifndef REPO_PATH
	$(error must set REPO_PATH, eg "make docker-tag-push REPO_PATH=${REPO}")
endif
ifndef PUBLISH
	docker tag ${IMG} ${REPO_PATH}/${IMG}
	docker push ${REPO_PATH}/${IMG}
else
	docker tag ${IMG} ${REPO_PATH}/${PUBLISH_IMG}
	docker push ${REPO_PATH}/${PUBLISH_IMG}
endif

varnish-controller: fmt vet
	go build -o ${ROOT_DIR}bin/varnish-controller ${ROOT_DIR}cmd/varnish-controller/

# Build the docker image with varnishd itself and varnish modules
docker-build-varnish:
	docker build --platform ${PLATFORM} ${ROOT_DIR} -t ${VARNISH_IMG} -f Dockerfile.varnishd

docker-tag-push-varnish:
ifndef REPO_PATH
	$(error must set REPO_PATH, eg "make docker-tag-push REPO_PATH=${REPO}")
endif
ifndef PUBLISH
	docker tag ${VARNISH_IMG} ${REPO_PATH}/${VARNISH_IMG}
	docker push ${REPO_PATH}/${VARNISH_IMG}
else
	docker tag ${VARNISH_IMG} ${REPO_PATH}/${VARNISH_PUBLISH_IMG}
	docker push ${REPO_PATH}/${VARNISH_PUBLISH_IMG}
endif

# Build the docker image with varnish controller
docker-build-varnish-controller: fmt vet
	docker build --platform ${PLATFORM} ${ROOT_DIR} -t ${VARNISH_CONTROLLER_IMG} -f Dockerfile.controller

docker-tag-push-varnish-controller:
ifndef REPO_PATH
	$(error must set REPO_PATH, eg "make docker-tag-push REPO_PATH=${REPO}")
endif
ifndef PUBLISH
	docker tag ${VARNISH_CONTROLLER_IMG} ${REPO_PATH}/${VARNISH_CONTROLLER_IMG}
	docker push ${REPO_PATH}/${VARNISH_CONTROLLER_IMG}
else
	docker tag ${VARNISH_CONTROLLER_IMG} ${REPO_PATH}/${VARNISH_CONTROLLER_PUBLISH_IMG}
	docker push ${REPO_PATH}/${VARNISH_CONTROLLER_PUBLISH_IMG}
endif

# Build the docker image with varnish metrics exporter
docker-build-varnish-exporter:
	docker build --platform ${PLATFORM} ${ROOT_DIR} -t ${VARNISH_METRICS_IMG} -f Dockerfile.exporter

docker-tag-push-varnish-exporter:
ifndef REPO_PATH
	$(error must set REPO_PATH, eg "make docker-tag-push REPO_PATH=${REPO}")
endif
ifndef PUBLISH
	docker tag ${VARNISH_METRICS_IMG} ${REPO_PATH}/${VARNISH_METRICS_IMG}
	docker push ${REPO_PATH}/${VARNISH_METRICS_IMG}
else
	docker tag ${VARNISH_METRICS_IMG} ${REPO_PATH}/${VARNISH_METRICS_PUBLISH_IMG}
	docker push ${REPO_PATH}/${VARNISH_METRICS_PUBLISH_IMG}
endif

docker-build-pod: docker-build-varnish docker-build-varnish-exporter docker-build-varnish-controller
docker-tag-push-pod: docker-tag-push-varnish docker-tag-push-varnish-exporter docker-tag-push-varnish-controller

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.11.3
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

e2e-tests:
	sh $(ROOT_DIR)hack/create_dev_cluster.sh
	KUBECONFIG=$(ROOT_DIR)e2e-tests-kubeconfig go test ./tests
	sh $(ROOT_DIR)hack/delete_dev_cluster.sh

KUSTOMIZE = $(shell pwd)/bin/kustomize
kustomize:
	$(call go-get-tool,$(KUSTOMIZE),sigs.k8s.io/kustomize/kustomize/v3@v3.8.7)

# Generate bundle manifests and metadata, then validate generated files.
.PHONY: bundle
bundle: manifests kustomize
	yq w -i config/manager/deployment.yaml 'spec.template.spec.containers(name==varnish-operator).env(name==CONTAINER_IMAGE).value' $(PUBLISH_IMG)
	yq w -i config/manifests/bases/varnish-operator.clusterserviceversion.yaml 'metadata.annotations.containerImage' $(PUBLISH_IMG)
	yq w -i config/manifests/bases/varnish-operator.clusterserviceversion.yaml 'metadata.annotations.createdAt' $(date +"%Y-%m-%d")
	operator-sdk generate kustomize manifests -q
	cd config/manager && $(KUSTOMIZE) edit set image controller=$(PUBLISH_IMG)
	$(KUSTOMIZE) build config/manifests | operator-sdk generate bundle -q --overwrite --version $(VERSION) $(BUNDLE_METADATA_OPTS)
	operator-sdk bundle validate ./bundle
	cp Dockerfile.bundle ./bundle/Dockerfile
	mv ./bundle ./$(VERSION)
