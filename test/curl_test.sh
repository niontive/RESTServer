#!/bin/bash
trap "kill 0" EXIT

# YAML files to test
declare -a yaml_files=("invalid1.yaml" \
                       "invalid2.yaml" \
                       "valid1.yaml" \
                       "valid2.yaml")

# Start RESTful server
start_server() {
    ../bin/server
}

# Test POST method
test_post() {
    curl -i -X POST localhost:10000/createmetadata \
         -H "Content-Type: text/x-yaml" \
         --data-binary @"$1"
}

start_server &
SERVER_PID=$!

for file in "${yaml_files[@]}"
do
    test_post ${file}
done

kill $SERVER_PID