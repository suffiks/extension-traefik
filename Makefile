generate:
	../../testerator/cmd/extgen/extgen crd -type Traefik
	../../testerator/cmd/extgen/extgen rbac -name traefik-extension

docker:
	docker build -t github.com/suffiks/extension-traefik:latest .

kind: docker
	kind load docker-image github.com/suffiks/extension-traefik:latest
	kubectl apply -k config
