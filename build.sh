#!/bin/sh

export DOCKER_BUILDKIT=1

docker build --target bin --output bin/ .