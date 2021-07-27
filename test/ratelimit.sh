#!/bin/bash

# Start RESTFul Server
start_server() {
    # TODO: make this configurable
    ../bin/server
}

# POST metadata
post_metadata() {
    curl -i -X POST localhost:10000/createmetadata \
         -H "Content-Type: text/x-yaml" \
         --data-binary @"$1"
}

start_server &
SERVER_PID=$!

for i in {0..100}
    do
        post_metadata "./yaml/valid1.yaml"

done

kill $SERVER_PID