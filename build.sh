#!/usr/bin/env sh

commit_id=$(git rev-parse --short HEAD)
image="largezhou/gin_starter"
image_with_tag="largezhou/gin_starter:$commit_id"

docker build . -t "$image_with_tag" -f ./Dockerfile

docker tag "$image_with_tag" "$image:latest"

docker push "$image_with_tag"
docker push "$image:latest"
