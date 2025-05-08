# Demo app

Demo application showcasing features like metrics exposition, ... .

## Deploy

```
TAG=1.0 make docker-build;
TAG=1.0 make docker-push;
```

## Load image to kind cluster

```
make kind-load
```