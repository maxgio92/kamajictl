# The [Kamaji](https://kamaji.clastix.io) CLI

## Installation

```shell
go install github.com/maxgio92/kamajictl@latest
```

## Getting started

### Install Kamaji

```shell
$ kamajictl install
 ►  applying resources
 ✔  resources applied
 ◎  verifying installation
 ✔  etcd: statefulset ready
 ✔  kamaji-controller-manager: deployment ready
 ✔  Install finished
```

### Create a `TenantControlPlane`

```shell
$ kamajictl create tcp foo -n default --service-type=ClusterIP --version=v1.26.1
✔  TenantControlPlane default/foo created
```

### List all your `TenantControlPlane`s

```shell
$ kamajictl get tcp
┌──────────────────────────────────────────────────────────────────────┐
| NAMESPACE | NAME | VERSION | STATUS | ENDPOINT          | DATA STORE |
| default   | bar  | v1.26.1 | Ready  | 10.96.219.55:6443 | default    |
| default   | foo  | v1.26.1 | Ready  | 10.96.1.206:6443  | default    |
└──────────────────────────────────────────────────────────────────────┘
```

### Get a `TenantControlPlane`'s `kubeconfig`

```shell
$ kamajictl get kubeconfig foo -n default -o yaml
apiVersion: v1
clusters:
- cluster:
    [...]
    server: https://10.96.1.206:6443
  name: foo
contexts:
- context:
    cluster: foo
    user: kubernetes-admin
  name: kubernetes-admin@foo
current-context: kubernetes-admin@foo
kind: Config
...
```

### Cleanup

```shell
$ kamajictl uninstall
 ►  Deleting Kamaji TCPs
 -  TCP bar/default deleted
 -  TCP foo/default deleted
 ✔  TenantControlPlanes deleted
 ►  Deleting Kamaji controller
 ✔  Kamaji controller uninstalled
 ►  Deleting Kamaji CRDs
 -  CRD datastores.kamaji.clastix.io deleted
 -  CRD tenantcontrolplanes.kamaji.clastix.io deleted
 ✔  Kamaji CRDs deleted
 ►  Deleting Kamaji webhook configurations
 -  validatingwebhookconfiguration kamaji-validating-webhook-configuration deleted
 -  mutatingwebhookconfiguration kamaji-mutating-webhook-configuration deleted
 ✔  Kamaji webhook configurations deleted
 ►  Deleting Kamaji RBAC
 -  clusterrole kamaji-manager-role deleted
 -  clusterrole kamaji-metrics-reader deleted
 -  clusterrole kamaji-proxy-role deleted
 -  clusterrolebinding kamaji-manager-rolebinding deleted
 -  clusterrolebinding kamaji-proxy-rolebinding deleted
 ✔  Kamaji RBAC deleted
 ✔  Uninstall finished
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

