## Development

### Prerequisites

##### Tools to build and test k8s-secret-injector

Apart from kubernetes cluster, there are some tools which are needed to build and test k8s-secret-injector.

- [Git](https://git-scm.com/downloads)
- [Go](https://golang.org/dl/)
- [Docker](https://docs.docker.com/install/)
- [Make](https://www.gnu.org/software/make/manual/make.html)

## Build Locally

To achieve this, execute this command:-

```shell
make build-code
```

## Build Image

k8s-secret-injector gets packaged as a container image for running on Kubernetes cluster. These instructions will guide you to build image.

```shell
make build-image
```

## Testing

For testing you may have to setup Hashicorp Vault.

## Run Tests
```shell
make test
```
