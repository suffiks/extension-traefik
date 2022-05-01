generate: extgen
	./bin/extgen crd -type Traefik
	./bin/extgen rbac -name traefik-extension

docker:
	docker build -t github.com/suffiks/extension-traefik:latest .

kind: docker
	kind load docker-image github.com/suffiks/extension-traefik:latest
	kubectl apply -k config


## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN): ## Ensure that the directory exists
	mkdir -p $(LOCALBIN)

## Tool Binaries
EXTGEN ?= $(LOCALBIN)/extgen

## Tool Versions
EXTGEN_VERSION ?= v0.1.0

.PHONY: extgen
extgen: $(EXTGEN)
$(EXTGEN): ## Download extgen locally if necessary.
	GOBIN=$(LOCALBIN) go install github.com/suffiks/suffiks/cmd/extgen@$(EXTGEN_VERSION)
