#!/bin/sh

export VERSION=$(git describe --tags)

if [ -z "$(docker buildx ls | grep 'linux/arm64')" ]; then 
    echo "Create a docker buildx"
    docker run --rm --privileged docker/binfmt:66f9012c56a8316f9244ffd7622d7c21c1f6f28d
    docker buildx create --name mybuilder  --platform linux/arm64
fi

docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 --push -t dil001/lms-control .
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 --push -t dil001/lms-control:$VERSION .