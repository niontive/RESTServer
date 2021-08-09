package main

import (
	"context"
	"errors"
	"net/http"

	"gopkg.in/yaml.v2"
)

//
// Retrieve metadata from metatdata store
//
func doGetAppMetaData(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	c := make(chan bool) // Channel for anonymous function
	go func(ctx context.Context) {
		logger.Info("EndpointHit: getAppMetaData")
		var metaDataFound bool = false

		// Retrieve and parse query string
		query := r.URL.Query()
		tmpMetaData := []appMetaData{}
		encoder := yaml.NewEncoder(w)
		dupTracker := make(map[string]bool)
		for k, v := range query {
			var errFunc error
			select {
			case <-ctx.Done(): // Context is done so don't perform more operations
				return
			default:
			}
			v = replaceUnderscore(k, v) // Underscore character may be treated as space
			logger.Infof("Query key: %v. Query value: %v.", k, v)

			// Perform metadata store search and output result as YAML format
			tmpMetaData, errFunc = dataStore.Search(k, v)
			if errFunc == nil {
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
				logger.Warnf("Unable to find metadata for key '%v': %v", k, errFunc)
			}
		}
		encoder.Close()
		c <- metaDataFound
	}(ctx)
	select {
	case <-ctx.Done(): // End operation if ctx.Done() occurs
		return ctx.Err()
	case result := <-c:
		if result == true {
			return nil
		} else {
			return errors.New("No metadata found")
		}
	}
}

//
// Handler for retrieving metadata entries
//
func getAppMetaData(w http.ResponseWriter, r *http.Request) {
	err := doGetAppMetaData(r.Context(), w, r)
	if err != nil {
		logger.Warnf("getAppMetaData error: %v", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}
