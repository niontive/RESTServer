package main

import (
	"context"
	"errors"
	"net/http"

	"gopkg.in/yaml.v2"
)

func dataStoreSearch(ctx context.Context, ch chan error, k string, v []string, md *[]appMetaData) {
	var err error
	logger.Info("Searching datastore")
	if ctx.Err() != nil {
		ch <- ctx.Err()
		return
	}
	*md, err = dataStore.Search(k, v)
	logger.Info("Finished adding to datastore")
	ch <- err
}

//
// Search the metadata store for an entry
//
func doDataStoreSearch(ctx context.Context, k string, v []string) (md []appMetaData, err error) {
	ch := make(chan error)

	go dataStoreSearch(ctx, ch, k, v, &md)

	select {
	case <-ctx.Done():
		logger.Warn("Request context is closed")
		return nil, ctx.Err()
	case err = <-ch:
		return
	}
}

//
// Retrieve metadata from metadata store
//
func doGetAppMetaData(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var metaDataFound bool = false

	logger.Info("EndpointHit: getAppMetaData")

	// Retrieve and parse query string
	query := r.URL.Query()
	tmpMetaData := []appMetaData{}
	encoder := yaml.NewEncoder(w)
	dupTracker := make(map[string]bool)

	for k, v := range query {
		var err error

		v = replaceUnderscore(k, v) // Underscore character may be treated as space
		logger.Infof("Query key: %v. Query value: %v.", k, v)

		// Perform metadata store search and output result as YAML format
		tmpMetaData, err = doDataStoreSearch(ctx, k, v)
		if err == nil {
			logger.Infof("Found metadata for key %v", k)
			for _, data := range tmpMetaData {
				// Check if we already retrieved this metadata
				// Since the database doesn't have duplicate titles,
				// we can use Title as a duplicate tracker
				if _, test := dupTracker[data.Title]; !test {
					dupTracker[data.Title] = true
					encoder.Encode(data)
					metaDataFound = true
				}
			}
		} else {
			logger.Warnf("Unable to find metadata for key '%v': %v", k, err)
		}
	}
	encoder.Close()

	if metaDataFound != true {
		logger.Warn("No metadata found")
		internalError := http.StatusInternalServerError
		http.Error(w, errors.New("No metadata found").Error(), internalError)
	}
}

//
// Handler for retrieving metadata entries
//
func getAppMetaData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	doGetAppMetaData(ctx, w, r)
}
