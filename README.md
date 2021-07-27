# Golang RESTful API Server for Application Metadata

## Overview
This codebase implements a RESTful API to store and retrieve application metadata.

## Preqrequisties
Below are preqrequisites. Versions listed have been validated. 
* Go - 1.13
* Docker - 20.10.7
* Curl - 7.68.0
* Ubuntu - 20.04.2 LTS

## Building Executable and Running Golang Unit Tests
Use script "build.sh" in the root directory. This script has three flags:
* -b : build the server. This uses Docker to build in a container. Output executable is located in the ./bin folder
* -c : clean. Removes ./bin directory and removes the Docker build cache
* -t : run unit tests. See file "server_test.go" for unit tests.

## Running the Server
After building, run the server by executing "./bin/server". The server will use port 10000 on your localhost.

To verify the server is running, open a web browser and go to "http://localhost:10000". You should see the
following text: "No application meta data stored!".

A POST request to endpoint "/createmetadata" will store metadata. Requests must:
* Be in YAML format. See "test/yaml/valid1.yaml" for which fields are supported.
* Contain nonempty values for all fields.
* Contain valid values for certain fields, including email and website.

A GET request to endpoint "/getmetadata" will retrieve metadata and return the YAML representaiton of the metadata. Use a query string to search for specific metadata. For example, 

    curl "http://localhost:10000/getmetadata?license=Apache-2.0&title=Valid_App_1"

would return all metadata that contains license "Apache-2.0" OR has title "Valid App 1". Do note that for certain fields,
underscores get replaced as spaces.

See the integration tests under /test for examples on performing these requests. Detailed information on what is/isn't supported and design decisions can be found in file /doc/features.md.