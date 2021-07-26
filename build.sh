#!/bin/sh

export DOCKER_BUILDKIT=1

build_server() {
    docker build --target bin --output bin/ .
}

run_unit_tests() {
    docker build --progress=plain --target unit-test .
}

clean() {
    rm -rf ./bin
    docker builder prune -f
}

while getopts "btc" flag
do
    case "${flag}" in
        b) build_server;;
        t) run_unit_tests;;
        c) clean;;
    esac
done