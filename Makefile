docker-build-sidecar-injector:
	docker build --no-cache -t ${CONTAINER_NAME}:${IMAGE_TAG} . -f dockerfiles/Dockerfile.sidecar-injector

docker-build-proxy:
	docker build --no-cache -t ${CONTAINER_NAME}:${IMAGE_TAG} . -f dockerfiles/Dockerfile.proxy

proxy:
	cd cmd/opentelemetry && go mod download && go mod tidy && go run main.go