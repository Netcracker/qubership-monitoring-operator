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
CRDS_HELM_CRDS_FOLDER := charts/qubership-monitoring-crds

# Directories and files
BUILD_DIR=build
OUTPUT_DIR=$(BUILD_DIR)/_output
CRDS_DIR=$(BUILD_DIR)/_crds

# CRDs inside the subcharts
MON_CRD_FOLDER=$(HELM_FOLDER)/crds
GRAFANA_CRD_FOLDER=$(HELM_FOLDER)/charts/grafana-operator/crds
PROM_OPER_CRD_FOLDER=$(HELM_FOLDER)/charts/prometheus-operator/crds
PROM_ADAPTER_CRD_FOLDER=$(HELM_FOLDER)/charts/prometheus-adapter-operator/crds
VM_CRD_FOLDER=$(HELM_FOLDER)/charts/victoriametrics-operator/crds

# Documents folders
DOC_FOLDER := docs
SITE_FOLDER := site

# Set build version
ARTIFACT_NAME="qubership-monitoring-operator"
VERSION?=0.88.0

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
DOCKERFILE=cmd/operator/Dockerfile

# CRD update tool and operator versions (override on the command line, e.g.
# `make update-prometheus-crds PROMETHEUS_OPERATOR_VERSION=0.92.1`)
CRD_UPDATE_TOOL=tools/crd-update/crd-update.py
PYTHON?=python3
# renovate: datasource=docker depName=quay.io/prometheus-operator/prometheus-operator
PROMETHEUS_OPERATOR_VERSION?=0.92.1
# renovate: datasource=docker depName=victoriametrics/operator
VICTORIAMETRICS_OPERATOR_VERSION?=0.73.1
GRAFANA_OPERATOR_VERSION?=5.24.0
LOCALBIN ?= $(CURDIR)/bin
CONTROLLER_GEN_VERSION ?= 0.21.0
CONTROLLER_GEN := $(LOCALBIN)/controller-gen-$(CONTROLLER_GEN_VERSION)
SETUP_ENVTEST_VERSION ?= v0.0.0-20250308055145-5fe7bb3edc86
SETUP_ENVTEST := $(LOCALBIN)/setup-envtest-$(SETUP_ENVTEST_VERSION)
ENVTEST_K8S_VERSION ?= 1.25.x!
YQ_VERSION ?= 4.53.2
YQ := $(LOCALBIN)/yq-$(YQ_VERSION)
YQ_OS := $(shell go env GOOS)
YQ_ARCH := $(shell go env GOARCH)
YQ_SHA256_linux_amd64 := d56bf5c6819e8e696340c312bd70f849dc1678a7cda9c2ad63eebd906371d56b
YQ_SHA256_linux_arm64 := 03061b2a50c7a498de2bbb92d7cb078ce433011f085a4994117c2726be4106ea
YQ_SHA256_darwin_amd64 := 616b0a0f6a5b79d746f05a169c2b9bb40dee00c605ef165b9a1c1681bba738ac
YQ_SHA256_darwin_arm64 := 541ba2287560df70f561955e2d7f7e1cd00cf2a15a884f6b5c87a4bfa887bc07
YQ_SHA256 := $(YQ_SHA256_$(YQ_OS)_$(YQ_ARCH))

CRD_COMPACTION_EXPRESSION := del(.. | select(tag == "!!map" and has("description") and (.description | tag == "!!str")).description) | del(.metadata.annotations."helm.sh/hook", .metadata.annotations."helm.sh/hook-weight")
PLATFORM_MONITORING_CRD := $(MON_CRD_FOLDER)/monitoring.netcracker.com_platformmonitorings.yaml
PLATFORM_MONITORING_COMPATIBILITY_EXPRESSION := del(.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.prometheus.properties.remoteRead.items.properties.url.pattern) | (.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.prometheus.properties.remoteRead.items.properties.oauth2.properties.tokenUrl |= (del(.pattern) | .minLength = 1)) | (.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.prometheus.properties.remoteWrite.items.properties.oauth2.properties.tokenUrl |= (del(.pattern) | .minLength = 1)) | .spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.victoriametrics.properties.vmAlertManager.properties.storage.properties.disableMountSubPath = {"type": "boolean"}

GENERATED_PATHS := \
	api/v1/zz_generated.deepcopy.go \
	$(MON_CRD_FOLDER) \
	$(GRAFANA_CRD_FOLDER) \
	$(PROM_OPER_CRD_FOLDER) \
	$(PROM_ADAPTER_CRD_FOLDER) \
	$(VM_CRD_FOLDER) \
	$(CRDS_HELM_CRDS_FOLDER)/crds

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
generate: controller-gen yq
	echo "=> Generate CRDs and deepcopy ..."
	$(CONTROLLER_GEN) crd:crdVersions={v1},maxDescLen=0 \
					  object:headerFile="tools/boilerplate.go.txt" \
					  paths="./..." \
					  output:crd:artifacts:config=charts/qubership-monitoring-operator/crds/
	"$(YQ)" eval -i '$(PLATFORM_MONITORING_COMPATIBILITY_EXPRESSION)' "$(PLATFORM_MONITORING_CRD)"

.PHONY: generate-all
generate-all:
	$(MAKE) generate
	$(MAKE) update-operators-crds
	$(MAKE) update-crds

.PHONY: verify-generated
verify-generated: generate-all
	$(MAKE) verify-generated-drift

.PHONY: verify-generated-drift
verify-generated-drift:
	@set -e; \
	if ! git diff --quiet -- $(GENERATED_PATHS); then \
	  echo "Generated files differ from the Git index:" >&2; \
	  git diff --name-status -- $(GENERATED_PATHS) >&2; \
	  exit 1; \
	fi; \
	untracked_files="$$(git ls-files --others --exclude-standard -- $(GENERATED_PATHS))"; \
	if [[ -n "$$untracked_files" ]]; then \
	  echo "Generated files are untracked:" >&2; \
	  printf '%s\n' "$$untracked_files" >&2; \
	  exit 1; \
	fi

.PHONY: controller-gen
controller-gen:
	@set -e; \
	mkdir -p "$(LOCALBIN)"; \
	if [[ ! -x "$(CONTROLLER_GEN)" ]]; then \
	  temporary_dir="$$(mktemp -d "$(LOCALBIN)/.controller-gen.XXXXXX")"; \
	  trap 'rm -rf "$$temporary_dir"' EXIT; \
	  GOBIN="$$temporary_dir" go install \
	    sigs.k8s.io/controller-tools/cmd/controller-gen@v$(CONTROLLER_GEN_VERSION); \
	  mv "$$temporary_dir/controller-gen" "$(CONTROLLER_GEN)"; \
	fi

#########
# Build #
#########

# Build manager binary
.PHONY: build-binary
build-binary: generate fmt vet
	echo "=> Build binary ..."
	$(GO_BUILD_RECIPE) -o bin/manager ./cmd/operator/

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

.PHONY: setup-envtest
setup-envtest:
	@set -e; \
	mkdir -p "$(LOCALBIN)"; \
	if [[ ! -x "$(SETUP_ENVTEST)" ]]; then \
	  temporary_dir="$$(mktemp -d "$(LOCALBIN)/.setup-envtest.XXXXXX")"; \
	  trap 'rm -rf "$$temporary_dir"' EXIT; \
	  GOBIN="$$temporary_dir" go install \
	    sigs.k8s.io/controller-runtime/tools/setup-envtest@$(SETUP_ENVTEST_VERSION); \
	  mv "$$temporary_dir/setup-envtest" "$(SETUP_ENVTEST)"; \
	fi

.PHONY: envtest
envtest: setup-envtest
	@set -e; \
	mkdir -p "$(LOCALBIN)/.cache" "$(LOCALBIN)/.data"; \
	kubebuilder_assets="$$( \
	  XDG_CACHE_HOME="$(LOCALBIN)/.cache" \
	  XDG_DATA_HOME="$(LOCALBIN)/.data" \
	  "$(SETUP_ENVTEST)" use -p path "$(ENVTEST_K8S_VERSION)" \
	)"; \
	XDG_CACHE_HOME="$(LOCALBIN)/.cache" \
	KUBEBUILDER_ASSETS="$$kubebuilder_assets" \
	go test -tags=envtest -race ./test/envtests/... -count=1

.PHONY: yq
yq:
	@set -e; \
	mkdir -p "$(LOCALBIN)"; \
	temporary_yq="$$(mktemp "$(LOCALBIN)/.yq-$(YQ_VERSION).XXXXXX")"; \
	trap 'rm -f "$$temporary_yq"' EXIT; \
	if [[ ! -x "$(YQ)" ]]; then \
	  if [[ -z "$(YQ_SHA256)" ]]; then \
	    echo "Unsupported yq platform: $(YQ_OS)/$(YQ_ARCH)" >&2; \
	    exit 1; \
	  fi; \
	  curl -sSLf -o "$$temporary_yq" \
	    "https://github.com/mikefarah/yq/releases/download/v$(YQ_VERSION)/yq_$(YQ_OS)_$(YQ_ARCH)"; \
	  if command -v sha256sum >/dev/null 2>&1; then \
	    actual_checksum="$$(sha256sum "$$temporary_yq" | awk '{print $$1}')"; \
	  else \
	    actual_checksum="$$(shasum -a 256 "$$temporary_yq" | awk '{print $$1}')"; \
	  fi; \
	  if [[ "$$actual_checksum" != "$(YQ_SHA256)" ]]; then \
	    echo "Checksum verification failed for yq $(YQ_VERSION) on $(YQ_OS)/$(YQ_ARCH)." >&2; \
	    exit 1; \
	  fi; \
	  chmod +x "$$temporary_yq"; \
	  mv "$$temporary_yq" "$(YQ)"; \
	fi

.PHONY: test-crd-compaction
test-crd-compaction: yq
	@set -e; \
	temp_dir=$$(mktemp -d); \
	trap 'rm -rf "$$temp_dir"' EXIT; \
	fixture="$$temp_dir/crd.yaml"; \
	second_fixture="$$temp_dir/second-crd.yaml"; \
	printf '%s\n' \
	  'apiVersion: apiextensions.k8s.io/v1' \
	  'kind: CustomResourceDefinition' \
	  'metadata:' \
	  '  annotations:' \
	  '    helm.sh/hook: crd-install' \
	  '    helm.sh/hook-weight: "-5"' \
	  '    example.com/retained: "true"' \
	  'spec:' \
	  '  versions:' \
	  '  - schema:' \
	  '      openAPIV3Schema:' \
	  '        description: root schema documentation' \
	  '        type: object' \
	  '        properties:' \
	  '          spec:' \
	  '            description: spec schema documentation' \
	  '            type: object' \
	  '            properties:' \
	  '              description:' \
	  '                description: user-provided description' \
	  '                type: string' \
	  '              properties:' \
	  '                description: user-provided properties' \
	  '                type: string' > "$$fixture"; \
	cp "$$fixture" "$$second_fixture"; \
	for crd in "$$temp_dir"/*.yaml; do \
	  "$(YQ)" eval -i '$(CRD_COMPACTION_EXPRESSION)' "$$crd"; \
	done; \
	first_hashes="$$(sha256sum "$$temp_dir"/*.yaml)"; \
	for crd in "$$temp_dir"/*.yaml; do \
	  "$(YQ)" eval -i '$(CRD_COMPACTION_EXPRESSION)' "$$crd"; \
	done; \
	second_hashes="$$(sha256sum "$$temp_dir"/*.yaml)"; \
	[[ "$$first_hashes" == "$$second_hashes" ]]; \
	"$(YQ)" eval -e '.spec.versions[0].schema.openAPIV3Schema | has("description") | not' "$$fixture" > /dev/null; \
	"$(YQ)" eval -e '.spec.versions[0].schema.openAPIV3Schema.properties.spec | has("description") | not' "$$fixture" > /dev/null; \
	"$(YQ)" eval -e '.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.description | has("description") | not' "$$fixture" > /dev/null; \
	"$(YQ)" eval -e '.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.description.type == "string"' "$$fixture" > /dev/null; \
	"$(YQ)" eval -e '.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties | has("properties")' "$$fixture" > /dev/null; \
	"$(YQ)" eval -e '.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.properties | has("description") | not' "$$fixture" > /dev/null; \
	"$(YQ)" eval -e '.spec.versions[0].schema.openAPIV3Schema.properties.spec.properties.properties.type == "string"' "$$fixture" > /dev/null; \
	"$(YQ)" eval -e '(.metadata.annotations | has("helm.sh/hook") or has("helm.sh/hook-weight")) | not' "$$fixture" > /dev/null; \
	"$(YQ)" eval -e '.metadata.annotations."example.com/retained" == "true"' "$$fixture" > /dev/null

#################
# Documentation #
#################

# Run document generation
.PHONY: docs
docs:
	@echo "=> 'make docs' is deprecated; use 'make update-crds' instead."
	$(MAKE) update-crds

##########################
# Update CRDs Helm chart #
##########################

# Copy CRDs from canonical chart directories to the qubership-monitoring-crds Helm chart
.PHONY: update-crds
update-crds:
	echo "=> Update CRDs in dedicated Helm chart ..."
	rm -f $(CRDS_HELM_CRDS_FOLDER)/crds/*.yaml $(CRDS_HELM_CRDS_FOLDER)/crds/*.yml
	cp $(MON_CRD_FOLDER)/* $(CRDS_HELM_CRDS_FOLDER)/crds/
	cp $(GRAFANA_CRD_FOLDER)/* $(CRDS_HELM_CRDS_FOLDER)/crds/
	cp $(PROM_OPER_CRD_FOLDER)/* $(CRDS_HELM_CRDS_FOLDER)/crds/
	cp $(PROM_ADAPTER_CRD_FOLDER)/* $(CRDS_HELM_CRDS_FOLDER)/crds/
	cp $(VM_CRD_FOLDER)/* $(CRDS_HELM_CRDS_FOLDER)/crds/

###############
# CRD update #
###############

# Download upstream CRDs for managed operators and write them into the
# corresponding chart `crds/` folders. Versions can be overridden via
# PROMETHEUS_OPERATOR_VERSION / VICTORIAMETRICS_OPERATOR_VERSION /
# GRAFANA_OPERATOR_VERSION variables.
.PHONY: update-operators-crds
update-operators-crds: update-prometheus-crds update-victoriametrics-crds update-grafana-crds compact-prometheus-adapter-crds

.PHONY: update-prometheus-crds
update-prometheus-crds:
	echo "=> Update prometheus-operator CRDs (v$(PROMETHEUS_OPERATOR_VERSION)) ..."
	$(PYTHON) $(CRD_UPDATE_TOOL) \
		--operator prometheus \
		--version $(PROMETHEUS_OPERATOR_VERSION) \
		--output-dir $(PROM_OPER_CRD_FOLDER)
	rm -f $(VM_CRD_FOLDER)/monitoring.coreos.com_*.yaml
	cp $(PROM_OPER_CRD_FOLDER)/monitoring.coreos.com_*.yaml $(VM_CRD_FOLDER)/

.PHONY: update-victoriametrics-crds
update-victoriametrics-crds:
	echo "=> Update victoriametrics-operator CRDs (v$(VICTORIAMETRICS_OPERATOR_VERSION)) ..."
	$(PYTHON) $(CRD_UPDATE_TOOL) \
		--operator victoriametrics \
		--version $(VICTORIAMETRICS_OPERATOR_VERSION) \
		--output-dir $(VM_CRD_FOLDER)

.PHONY: update-grafana-crds
update-grafana-crds: yq
	echo "=> Update grafana-operator CRDs (v$(GRAFANA_OPERATOR_VERSION)) ..."
	$(PYTHON) $(CRD_UPDATE_TOOL) \
		--operator grafana \
		--version $(GRAFANA_OPERATOR_VERSION) \
		--output-dir $(GRAFANA_CRD_FOLDER)
	for crd in "$(GRAFANA_CRD_FOLDER)"/*.yaml; do \
	  "$(YQ)" eval -i '$(CRD_COMPACTION_EXPRESSION)' "$$crd"; \
	done

.PHONY: compact-prometheus-adapter-crds
compact-prometheus-adapter-crds: yq
	echo "=> Compact prometheus-adapter-operator CRDs ..."
	for crd in "$(PROM_ADAPTER_CRD_FOLDER)"/*.yaml; do \
	  "$(YQ)" eval -i '$(CRD_COMPACTION_EXPRESSION)' "$$crd"; \
	done

###################
# Running locally #
###################

# Run against the configured Kubernetes cluster in ~/.kube/config
.PHONY: run
run: generate fmt vet
	echo "=> Run ..."
	go run ./cmd/operator/

#################
# Building docs #
#################

# Install the dependencies
.PHONY: install-site-dependencies
install-site-dependencies:
	echo "=> Install site dependencies ..."
	pip install -r site/requirements.txt

# Prepare the docs directory
.PHONY: prepare-site-directory
prepare-site-directory:
	echo "=> Prepare site directory ..."
	rm -rf $(SITE_FOLDER)/docs
	mkdir -p $(SITE_FOLDER)
	cp -rL $(DOC_FOLDER) $(SITE_FOLDER)/

# Build the docs
.PHONY: build-site
build-site: prepare-site-directory install-site-dependencies
	echo "=> Build site ..."
	zensical build -f site/mkdocs.yml --clean


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

	# Create a directory for the flat CRD archive
	mkdir -p "$(CRDS_DIR)"

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

	# The standalone CRD chart is the generated, deduplicated aggregate.
	cp charts/qubership-monitoring-crds/crds/* "${BUILD_DIR}/_crds/"

	# Navigate to dir to avoid unnecessary directories in result archive\
	# name like: qubership-monitoring-operator-0.60.0-crds.zip
	cd "${BUILD_DIR}"/_crds && zip -r "../../${OUTPUT_DIR}/${ARTIFACT_NAME}-${VERSION}-crds.zip" ./*
