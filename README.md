## Tech stack
- Backend building blocks:

    - [gofiber/fiber/v2](https://github.com/gofiber/fiber)
    - [spf13/cobra](https://github.com/spf13/cobra)

- Utils:

    - [https://github.com/ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv)
    - k8s.io:

        - k8s.io/client-go/kubernetes
        - k8s.io/api
        - k8s.io/apimachinery
    - [open-telemetry/opentelemetry-go](https://github.com/open-telemetry/opentelemetry-go)
    - github.com/gofiber/contrib/otelfiber/v2
    

## Proxy
- Copy config:
```bash
cp cmd/opentelemetry/config.example.yaml cmd/opentelemetry/config.yaml
```

1. Start service
```bash
cd cmd/opentelemetry && go mod download && go mod tidy && go run main.go
or
make proxy
```

2. Build image:
```bash
export CONTAINER_NAME=your_container_name
export IMAGE_TAG=your_image_tag
```

```bash
make docker-build-proxy
or
docker build --no-cache -t ${CONTAINER_NAME}:${IMAGE_TAG} . -f dockerfiles/Dockerfile.proxy
```


## Sidecar
#### Coding sample: [kubernetes-sidecar-injector](https://github.com/ExpediaGroup/kubernetes-sidecar-injector)

1. Build image:
```bash
export CONTAINER_NAME=your_container_name
export IMAGE_TAG=your_image_tag
```

```bash
make docker-build-sidecar-injector
or
docker build --no-cache -t ${CONTAINER_NAME}:${IMAGE_TAG} . -f dockerfiles/Dockerfile.sidecar-injector
```

2. Release sidecar-injector:

- Deploy with helm
```bash
## install helm
brew install helm
```

```bash
## helm install chart
helm upgrade -i sidecar-injector ./deployment/charts/sidecar-injector/. \
--namespace=sidecar-injector --create-namespace \
--set image.repository=${CONTAINER_NAME} \
--set image.tag=${IMAGE_TAG}
```

- Deploy with terraform
```bash
## install terraform
brew install terraform
```

```bash
cd deployment/terraform
terraform init
terraform plan -target=module.sidecar-injector
terraform apply -target=module.sidecar-injector -auto-approve
```