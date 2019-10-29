# Image URL to use in all building/pushing image targets
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION ?= $(shell cat ${ROOT_DIR}version.txt)
PUBLISH_IMG ?= varnish-operator:${VERSION}
VARNISH_PUBLISH_IMG ?= varnish:${VERSION}
VARNISH_IMG ?= ${VARNISH_PUBLISH_IMG}-dev
IMG ?= ${PUBLISH_IMG}-dev
NAMESPACE ?= "default"
CRD_OPTIONS ?= "crd:trivialVersions=true"

# all: test manager
all: test manager varnish-controller

# Run tests
test: generate fmt vet manifests
	go test icm-varnish-k8s-operator/pkg/... icm-varnish-k8s-operator/cmd/... icm-varnish-k8s-operator/api/... -coverprofile cover.out

# Run lint tools
lint:
	golangci-lint run

# Build manager binary
manager: generate fmt vet
	go build -o ${ROOT_DIR}bin/manager icm-varnish-k8s-operator/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	NAMESPACE=${NAMESPACE} LOGLEVEL=debug LOGFORMAT=console CONTAINER_IMAGE=us.icr.io/icm-varnish/${VARNISH_IMG} LEADERELECTION_ENABLED=false WEBHOOKS_ENABLED=false go run ${ROOT_DIR}cmd/manager/main.go

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
install: manifests
	kustomize build ${ROOT_DIR}config/crd | kubectl apply -f -

uninstall:
	kustomize build ${ROOT_DIR}config/crd | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=varnish-operator paths="./..." output:crd:artifacts:config=config/crd/bases
	kustomize build ${ROOT_DIR}config/crd > $(ROOT_DIR)varnish-operator/templates/manager_customresourcedefinition.yaml
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=varnish-operator paths="./..." output:crd:none output:rbac:stdout > $(ROOT_DIR)varnish-operator/templates/manager_clusterrole.yaml

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
	docker build ${ROOT_DIR} -t ${IMG} -f Dockerfile

# Tag and push the docker image
docker-tag-push:
ifndef REPO_PATH
	$(error must set REPO_PATH, eg "make docker-tag-push REPO_PATH=us.icr.io/icm-varnish")
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

# Build the docker image
docker-build-varnish: fmt vet
	docker build ${ROOT_DIR} -t ${VARNISH_IMG} -f Dockerfile.Varnish

docker-tag-push-varnish:
ifndef REPO_PATH
	$(error must set REPO_PATH, eg "make docker-tag-push REPO_PATH=us.icr.io/icm-varnish")
endif
ifndef PUBLISH
	docker tag ${VARNISH_IMG} ${REPO_PATH}/${VARNISH_IMG}
	docker push ${REPO_PATH}/${VARNISH_IMG}
else
	docker tag ${VARNISH_IMG} ${REPO_PATH}/${VARNISH_PUBLISH_IMG}
	docker push ${REPO_PATH}/${VARNISH_PUBLISH_IMG}
endif

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.0
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
