# The [Kamaji](https://kamaji.clastix.io) CLI

## Installation

```shell
go install github.com/maxgio92/kamajicli@latest
```

## Getting started

### Create a `TenantControlPlane`

```shell
kamaji create tcp NAME -n NAMESPACE [--service-type=<SERVICE_TYPE>] [--version=<KUBERNETES_VERSION>]
```

### Get a `TenantControlPlane`'s `kubeconfig`

```shell
kamaji get kubeconfig NAME -n NAMESPACE
```

For other commands please read the [documentation](./docs).

## Development

### Lint

```shell
make lint
```

### Build

```shell
make build
```

### Run tests

```shell
make test
```

### Update the documentation

```shell
make docs
```

