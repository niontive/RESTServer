# Features

## Overview
This document outlines the interesting features, limitations and use of this RESTful API server.

## Design Decisions
Design decisions were made to improve quality and reliability. This includes:
* RWMutex for concurrent access. Since the metadata store resides in memory, we must be careful of concurrent access; for example, a read request may occur while a write request is being performed. RWMuticies protect all read/write operations for the global metadata store. This include the slice that holds the metadata as well as the map that holds duplicate information.
* Rate limiter for HTTP server. Spamming POST/GET requests may cause resource starvation, thereby affecting API availability. The rate limiter limits the number of requests to the HTTP server. This is set at one request for every 25 milliseconds, and based on hardware/software configuration, this rate may be increased or decreased.
* Timeout handler for POST and GET APIs. To prevent clients from waiting indefinitely, the requests should timeout after a certain time period. The timeout is set a five seconds, which may be increased or decreased.
* Cancellation signals via golang contexts. The HTTP request context is passed between APIs. If an API is slow, the context's done channel may trigger before API completion. This aids with API cancellation due to a timeout or a user cancellation. 

## Information on Usage
Below are important usage tips. This includes what is/isn't supported by the server. 
* YAML payloads must follow the guidelines set in APIServerExercise.md. 
* A payload is invalid if any field is empty.
* The create metadata API only supports storing one YAML file at a time
* When the create metadata API is called, validation is performed. Besides empty string checking, certain fields have additional validaiton. The email field of "maintainers" must meet the RFC5322 standard. All URLs must pass the golang ParseRequestURI() function.
* Each metadata must have a unique "title" field. When someone attempts to add a metadata to the store, the server will check if a metadata with the same "title" already exists; if so, the new metadata will not be added. One future improvement is tieing uniqueness to title AND app version.
* Requests to GET metadata must use a query string. An example valid request is:

    curl "http://localhost:10000/getmetadata?license=Apache-2.0&title=Valid_App_1"

Here, the text after the "?" sign are the key value pairs to be searched. All keys must be lowercase. The server will respond with the YAML representation of metadata the matches EITHER of the key/value pairs. 

To support spaces, an underscore may be used. Underscores are replaced with spaces for the following fields: title, name, company, license and description. For these fields, underscores are not permitted when the metadata is POSTed.
* The GET request does not send the same metadata twice, even if the metadata matches more than one key/value pair requested.
* If using curl, remember to use quotes around the URL for GET requests.