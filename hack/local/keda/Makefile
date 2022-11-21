REGISTRY_PORT ?= 5001
REGISTRY_NAME ?= kyma-registry.localhost
REGISTRY_CERT ?= registry.pem

CLUSTER_NAME ?= keda-local-dev

MANAGER_IMAGE_NAME ?= keda-manager-dev-local
MANAGER_IMAGE_TAG ?= 0.0.1

PROJECT_ROOT ?= ../../..

TEMPLATE_PATH ?= ${PROJECT_ROOT}/template.yaml

KYMA ?= ${PROJECT_ROOT}/bin/kyma-unstable

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
	@make -C ${PROJECT_ROOT} module-image \
		IMG=localhost:${REGISTRY_PORT}/${MANAGER_IMAGE_NAME}:${MANAGER_IMAGE_TAG}

.PHONY: module-build
module-build:
	@make -C ${PROJECT_ROOT} module-build \
		IMG=k3d-${REGISTRY_NAME}:5000/${MANAGER_IMAGE_NAME}:${MANAGER_IMAGE_TAG} \
		MODULE_REGISTRY=localhost:${REGISTRY_PORT}/unsigned

.PHONY: fix-template
fix-template:
	@cat ${PROJECT_ROOT}/template.yaml \
	| sed -e 's/remote/control-plane/g' \
		-e 's/${REGISTRY_PORT}/5000/g'  \
	      	-e 's/localhost/k3d-${REGISTRY_NAME}/g' \
	> ${PROJECT_ROOT}/template-k3d.yaml

.PHONY: enable-module
enable-module:
	kubectl patch kymas.operator.kyma-project.io \
		default-kyma \
		-n kcp-system \
		--type=merge \
		--patch-file ./k3d-patches/patch-kyma.yaml

.PHONY: patch-mod-mgr-role
patch-mod-mgr-role:
	kubectl patch clusterrole \
		module-manager-manager-role \
		--patch-file=./k3d-patches/patch-k3d-module-manager-clusterrole.yaml

.PHONY: patch-coredns
patch-coredns:
	kubectl patch -n kube-system cm coredns --patch-file ./k3d-patches/patch-k3d-coredns.yaml

.PHONY: install-module
install-module:
	${KYMA} alpha deploy \
		--template=${PROJECT_ROOT}/template-k3d.yaml

.PHONY: run
run: create-k3d patch-coredns module-image module-build fix-template install-module patch-mod-mgr-role enable-module

.PHONY: stop
stop: delete-k3d

.PHONY: verify
verify:
	@${PROJECT_ROOT}/hack/verify_kyma_status.sh

