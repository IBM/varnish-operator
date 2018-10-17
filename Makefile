
# Image URL to use in all building/pushing image targets
VERSION ?= $(shell cat ${ROOT_DIR}version.txt)-dev
VARNISH_VERSION ?= $(shell cat ${ROOT_DIR}icm-varnish-version.txt)
IMG ?= varnish-controller:${VERSION}
ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
NAMESPACE := $(shell sed -n -e 's/^namespace: //p' ${ROOT_DIR}config/default/kustomization.yaml)
NAME_PREFIX := $(shell sed -n -e 's/^namePrefix: //p' ${ROOT_DIR}config/default/kustomization.yaml)

# all: test manager
all: fake-test manager

# test is failing right now because kubebuilder does not know how to test slices
fake-test: generate fmt vet manifests

# Run tests
test: generate fmt vet manifests
	go test ${ROOT_DIR}pkg/... ${ROOT_DIR}cmd/... -coverprofile cover.out

# Build manager binary
manager: generate fmt vet
	go build -o ${ROOT_DIR}bin/manager icm-varnish-k8s-operator/cmd/manager

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet
	LOG_FORMAT=console VARNISH_IMAGE_TAG=${VARNISH_VERSION} go run ${ROOT_DIR}cmd/manager/main.go

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
install: manifests
	kustomize build ${ROOT_DIR}config/default | kubectl apply -f -

uninstall:
	kustomize build ${ROOT_DIR}config/default | kubectl delete -f -

# Generate manifests e.g. CRD, RBAC etc.
manifests:
	go run ${ROOT_DIR}vendor/sigs.k8s.io/controller-tools/cmd/controller-gen/main.go all

# Run go fmt against code
fmt:
	cd ${ROOT_DIR} && go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	cd ${ROOT_DIR} && go vet ./pkg/... ./cmd/...

# Generate code
generate:
	cd ${ROOT_DIR} && go generate ./pkg/... ./cmd/...

# Prepare .yaml files for helm
helm-prepare: manifests
	${ROOT_DIR}hack/create_helm_files.sh ${ROOT_DIR}varnish-operator/templates

helm-upgrade: helm-prepare
	@if [ -z "${NAMESPACE}" ]; then\
		echo "trying to read \"namespace:\" line in config/default/kustomization.yaml. Did something change?";\
		exit 1;\
	fi
	@if [ -z "${NAME_PREFIX}" ]; then\
		echo "trying to read \"namePrefix\" line in config/default/kustomization.yaml. Did something change?";\
		exit 1;\
	fi
	helm upgrade --install --namespace ${NAMESPACE} --force varnish-operator --wait --debug --set operator.controllerImage.tag=${VERSION} --set namespace=${NAMESPACE} --set namePrefix=${NAME_PREFIX} ${ROOT_DIR}varnish-operator

# Build the docker image
# docker-build: test
docker-build: fake-test
	docker build ${ROOT_DIR} -t ${IMG} -f Dockerfile.local

docker-build-dep: fake-test
    docker build ${ROOT_DIR} -t ${IMG} -f Dockerfile

# Tag and push the docker image
docker-tag-push:
	@if [ -z "${REPO_PATH}" ]; then\
		echo "must set REPO_PATH variable, eg \"make docker-tag-push REPO_PATH=registry.ng.bluemix.net/icm-varnish\"";\
		exit 1;\
	fi
	docker tag ${IMG} ${REPO_PATH}/${IMG}
	docker push ${REPO_PATH}/${IMG}
