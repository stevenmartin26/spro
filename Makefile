

.PHONY: clean all init generate generate_mocks

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: generate
	go mod tidy
	go mod vendor

test:
	go test -short -coverprofile coverage.out -v ./...

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

generate_mocks: generate_repo_mocks generate_service_mocks

INTERFACES_REPO_GO_FILES := $(shell find repository -name "interfaces.go")
INTERFACES_REPO_GEN_GO_FILES := $(INTERFACES_REPO_GO_FILES:%.go=%.mock.gen.go)

generate_repo_mocks: $(INTERFACES_REPO_GEN_GO_FILES)
$(INTERFACES_REPO_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))

INTERFACES_SVC_GO_FILES := $(shell find service -name "interfaces.go")
INTERFACES_SVC_GEN_GO_FILES := $(INTERFACES_SVC_GO_FILES:%.go=%.mock.gen.go)

generate_service_mocks: $(INTERFACES_SVC_GEN_GO_FILES)
$(INTERFACES_SVC_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))