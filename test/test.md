# Running the Integration Tests

## Overview
There are three integration test scripts. Before running, build the server using the instructions found in the README.

YAML files used are in the /yaml directory

## Tests
* post.sh: test the "createmetadata" API. This script attempts to send two invalid and two valid YAML files to the server. View the log output for test results. There should be errors associated with POSTing the invalid YAML files, and the valid YAML files should be stored successfully.

* get.sh: test the "getmetadata" API. This script sends two YAML files to the server. The script then queries the server for each metadata via the application titles. The log output should show two retrieved YAML files. 

* ratelimit.sh: test the ratelimit feature of the server. The script spams the createmetadata API. After execution, the log output should show the server refusing requests because too many requests are made.