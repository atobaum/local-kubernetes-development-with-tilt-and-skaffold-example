An example of a local kubernetes development environment using [k3d](https://k3d.io/), [skaffold](https://skaffold.dev/), and [tilt](https://tilt.dev/).

# Prerequisites
```shell
brew install k3d skaffold tilt
```

# How to use
## Create k3d cluster
```shell
k3d cluster create local-dev --registry-create local-dev-registry
```

## 1. Using TILT
```shell
tilt up
```

## 2. Using SKAFFOLD
```shell
skaffold dev --port-forward 
```
