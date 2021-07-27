#!/bin/bash

# YAML files to test
declare -a yaml_files=("./yaml/valid1.yaml" \
                       "./yaml/valid2.yaml")

declare -a app_titles=("Valid_App_1" \
                       "Valid_App_2")

# Start RESTFul Server
start_server() {
    # TODO: make this configurable
    ../bin/server
}

# Test GET method
test_get() {
    curl http://localhost:10000/getmetadata?title=${1}
}

# POST metadata
post_metadata() {
    curl -i -X POST localhost:10000/createmetadata \
         -H "Content-Type: text/x-yaml" \
         --data-binary @"$1"
}

trap "exit" INT TERM ERR
trap "kill 0" EXIT

start_server &
sleep 1

for file in "${yaml_files[@]}"
do
    post_metadata ${file}
    sleep 1
done

for title in "${app_titles[@]}"
do
    test_get ${title}
    sleep 1
done
