SHELL=/usr/bin/env bash -o pipefail

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
	GOBIN=$(shell go env GOPATH)/bin
else
	GOBIN=$(shell go env GOBIN)
endif

#############
# Constants #
#############

# Helm charts directory
HELM_FOLDER := charts/qubership-monitoring-operator

# Directories and files
BUILD_DIR=build
OUTPUT_DIR=$(BUILD_DIR)/_output
CRDS_DIR=$(BUILD_DIR)/_crds

# CRDs inside the subcharts
MON_CRD_FOLDER=$(HELM_FOLDER)/crds
GRAFANA_CRD_FODLER=$(HELM_FOLDER)/charts/grafana-operator/crds
PROM_OPER_CRD_FOLDER=$(HELM_FOLDER)/charts/prometheus-operator/crds
PROM_ADAPTER_CRD_FOLDER=$(HELM_FOLDER)/charts/prometheus-adapter-operator/crds
VM_CRD_FOLDER=$(HELM_FOLDER)/charts/victoriametrics-operator/crds

# Documents folders
DOC_FOLDER := docs
CRD_DOC_FOLDER=$(DOC_FOLDER)/crds

# Set build version
ARTIFACT_NAME="qubership-monitoring-operator"
VERSION?=0.75.0

# Detect the build environment, local or Jenkins builder
BUILD_DATE=$(shell date +"%Y%m%d-%T")
ifndef JENKINS_URL
	BUILD_USER?=$(USER)
	BUILD_BRANCH?=$(shell git branch --show-current)
	BUILD_REVISION?=$(shell git rev-parse --short HEAD)
else
	BUILD_USER=$(BUILD_USER)
	BUILD_BRANCH=$(LOCATION:refs/heads/%=%)
	BUILD_REVISION=$(REPO_HASH)
endif

# The Prometheus common library import path
MONITORING_OPERATOR_PKG=github.com/Netcracker/qubership-monitoring-operator

# The ldflags for the go build process to set the version related data.
GO_BUILD_LDFLAGS=\
	-s \
	-X $(MONITORING_OPERATOR_PKG)/version.Revision=$(BUILD_REVISION) \
	-X $(MONITORING_OPERATOR_PKG)/version.BuildUser=$(BUILD_USER) \
	-X $(MONITORING_OPERATOR_PKG)/version.BuildDate=$(BUILD_DATE) \
	-X $(MONITORING_OPERATOR_PKG)/version.Branch=$(BUILD_BRANCH) \
	-X $(MONITORING_OPERATOR_PKG)/version.Version=$(VERSION)

# Go build flags
GO_BUILD_RECIPE=\
	GOOS=$(GOOS) \
	GOARCH=$(GOARCH) \
	CGO_ENABLED=0 \
	go build -ldflags="$(GO_BUILD_LDFLAGS)"

# Default test arguments
TEST_RUN_ARGS=-vet=off --shuffle=on

# List of packages, exclude integration tests that require "envtest"
pkgs = $(shell go list ./... | grep -v /test/envtests)
#pkgs += $(shell go list $(MONITORING_OPERATOR_PKG)/api...)

# Container name
CONTAINER_CLI?=docker
CONTAINER_NAME="qubership-monitoring-operator"
DOCKERFILE=Dockerfile

###########
# Generic #
###########

# Default run without arguments
.PHONY: all
all: generate test build-binary image docs archives

# Run only build
.PHONY: build
build: generate build-binary image docs archives

# Run only build inside the Dockerfile
.PHONY: build-image
build-image: generate image docs archives

# Remove all files and directories ignored by git
.PHONY: clean
clean:
	echo "=> Cleanup repository ..."
	git clean -Xfd .

##############
# Generating #
##############

# Generate code
.PHONY: generate
generate: controller-gen
	echo "=> Generate CRDs and deepcopy ..."
	$(CONTROLLER_GEN) crd:crdVersions={v1},maxDescLen=256 \
					  object:headerFile="tools/boilerplate.go.txt" \
					  paths="./..." \
					  output:crd:artifacts:config=charts/qubership-monitoring-operator/crds/

	if [[ "$$OSTYPE" == "darwin"* ]]; then \
	  SED_CMD="sed -i '' -e"; \
	else \
	  SED_CMD="sed -i"; \
	fi; \
	find charts/qubership-monitoring-operator/crds -name '*.yaml' | while read f; do \
	  $$SED_CMD "/^    controller-gen.kubebuilder.io.version.*/a\\    helm.sh/hook-weight: \"-5\"" "$$f"; \
	  $$SED_CMD "/^    controller-gen.kubebuilder.io.version.*/a\\    helm.sh/hook: crd-install" "$$f"; \
	done

# Find or download controller-gen
# download controller-gen if necessary
.PHONY: controller-gen
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.15.0 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif

#########
# Build #
#########

# Build manager binary
.PHONY: build-binary
build-binary: generate fmt vet
	echo "=> Build binary ..."
	$(GO_BUILD_RECIPE) -o bin/manager main.go

# Run go fmt against code
.PHONY: fmt
fmt:
	echo "=> Formatting Golang code ..."
	go fmt ./...

# Run go vet against code
.PHONY: vet
vet:
	echo "=> Examines Golang code ..."
	go vet ./...

###############
# Build image #
###############

.PHONY: image
image:
	echo "=> Build image ..."
	docker build --pull -t $(CONTAINER_NAME) -f $(DOCKERFILE) .

	# Set image tag if build inside the Jenkins
	for id in $(DOCKER_NAMES) ; do \
		docker tag $(CONTAINER_NAME) "$$id"; \
	done

###########
# Testing #
###########

.PHONY: test
test: unit-test

# Run unit tests in all packages
.PHONY: unit-test
unit-test:
	echo "=> Run Golang unit-tests ..."
	go test -race $(TEST_RUN_ARGS) $(pkgs) -count=1 -v

#################
# Documentation #
#################

# Run document generation
.PHONY: docs
docs: docs/crd/v1

.PHONY: docs/crd/v1
docs/crd/v1:
	echo "=> Copy CRDs from charts to documentation ..."
	rm -rf $(CRD_DOC_FOLDER)/*.yaml
	cp $(MON_CRD_FOLDER)/* $(CRD_DOC_FOLDER)/
	cp $(GRAFANA_CRD_FODLER)/* $(CRD_DOC_FOLDER)/
	cp $(PROM_OPER_CRD_FOLDER)/* $(CRD_DOC_FOLDER)/
	cp $(PROM_ADAPTER_CRD_FOLDER)/* $(CRD_DOC_FOLDER)/
	cp $(VM_CRD_FOLDER)/* $(CRD_DOC_FOLDER)/

###################
# Running locally #
###################

# Run against the configured Kubernetes cluster in ~/.kube/config
.PHONY: run
run: generate fmt vet
	echo "=> Run ..."
	go run ./main.go

############
# Archives #
############

# Run archives with helm chart and crds creation
.PHONY: archives
archives: cleanup prepare-charts archive-helm-chart archive-crds

# Remove build dir
.PHONY: cleanup
cleanup:
	rm -rf $(BUILD_DIR)

# Copy Helm charts to the /helm directory because the builder expect it in this dir
.PHONY: prepare-charts
prepare-charts:
	echo "=> Copy Helm charts to contract directory for build ..."
	mkdir -p $(OUTPUT_DIR)

	# Create directories to copy CRDs
	mkdir -p "$(CRDS_DIR)/qubership-monitoring-operator" \
	"$(CRDS_DIR)/prometheus-adapter-operator" \
	"$(CRDS_DIR)/prometheus-operator" \
	"$(CRDS_DIR)/victoriametrics-operator" \
	"$(CRDS_DIR)/grafana-operator"

# Archive Helm chart separately from application manifest
.PHONY: archive-helm-chart
archive-helm-chart:
	echo "=> Archive Helm charts ..."

	# Navigate to dir to avoid unnecessary directories in result archive
	# name like: qubership-monitoring-operator-0.60.0-chart.zip
	cd ./charts && zip -r "../${OUTPUT_DIR}/${ARTIFACT_NAME}-${VERSION}-chart.zip" ./*

# Archive CRDs separately from helm chart
.PHONY: archive-crds
archive-crds:
	echo "=> Archive CRDs ..."
	# Copy documentation how to apply CRDS
	cp docs/user-guides/manual-create-crds.md "${BUILD_DIR}"/_crds/README.md

	# Copy CRDs from different places in helm chart and subcharts
	cp charts/qubership-monitoring-operator/crds/* "${BUILD_DIR}/_crds/qubership-monitoring-operator/"
	cp charts/qubership-monitoring-operator/charts/prometheus-adapter-operator/crds/* "${BUILD_DIR}/_crds/prometheus-adapter-operator/"
	cp charts/qubership-monitoring-operator/charts/prometheus-operator/crds/* "${BUILD_DIR}/_crds/prometheus-operator/"
	cp charts/qubership-monitoring-operator/charts/victoriametrics-operator/crds/* "${BUILD_DIR}/_crds/victoriametrics-operator/"
	cp charts/qubership-monitoring-operator/charts/grafana-operator/crds/* "${BUILD_DIR}/_crds/grafana-operator/"

	# Navigate to dir to avoid unnecessary directories in result archive\
	# name like: qubership-monitoring-operator-0.60.0-crds.zip
	cd "${BUILD_DIR}"/_crds && zip -r "../../${OUTPUT_DIR}/${ARTIFACT_NAME}-${VERSION}-crds.zip" ./*
