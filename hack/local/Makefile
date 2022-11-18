REGISTRY_PORT ?= 5001
REGISTRY_NAME ?= kyma-registry
REGISTRY_CERT ?= registry.pem

CLUSTER_NAME ?= keda-local-dev.localhost

MANAGER_IMAGE_NAME ?= keda-manager-dev-local
MANAGER_IMAGE_TAG ?= 0.0.1

TEMPLATE_PATH ?= ../../template.yaml

KYMA ?= ../../bin/kyma-unstable

.PHONY: create-k3d-registry
create-k3d-registry:
	k3d registry create ${REGISTRY_NAME} --port 0.0.0.0:${REGISTRY_PORT} --no-help

.PHONY: delete-k3d-registry
delete-k3d-registry:
	-k3d registry delete ${REGISTRY_NAME}

.PHONY: create-k3d-cluster
create-k3d-cluster:
	k3d cluster create ${CLUSTER_NAME} \
	--registry-use=k3d-${REGISTRY_NAME}:${REGISTRY_PORT}

.PHONY: delete-k3d-cluster
delete-k3d-cluster:
	-k3d cluster delete ${CLUSTER_NAME}

.PHONY: create-k3d 
create-k3d: create-k3d-registry create-k3d-cluster

.PHONY: delete-k3d
delete-k3d: delete-k3d-cluster delete-k3d-registry

.PHONY: module-image
module-image:
	make -C ../.. module-image IMG=localhost:${REGISTRY_PORT}/${MANAGER_IMAGE_NAME}:${MANAGER_IMAGE_TAG}

.PHONY: module-build
module-build:
	make -C ../.. module-build \
		IMG=${REGISTRY_NAME}:${REGISTRY_PORT}/${MANAGER_IMAGE_NAME}/${MANAGER_IMAGE_TAG} \
		MODULE_REGISTRY=localhost:${REGISTRY_PORT}/unsigned

.PHONY: fix-template
fix-template:
	cat ../../template.yaml | sed -e 's/remote/control-plane/g' -e 's/${REGISTRY_PORT}/5000/g'  -e 's/localhost/k3d-${REGISTRY_NAME}/g'> ../../template-k3d.yaml

.PHONY: install-module
install-module:
	${KYMA} alpha deploy --template=../../template-k3d.yaml

.PHONY: install-dev
install-dev: create-k3d module-image module-build fix-template install-module


