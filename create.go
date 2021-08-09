package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

//
// Add new metadata entry to metadata store
//
func doCreateNewAppMetaData(ctx context.Context, r *http.Request) error {
	c := make(chan error) // Channel for anonymous function
	go func(ctx context.Context) {
		logger.Info("Endpoint Hit: createNewAppMetaData")

		// Read request into metadata structure
		reqBody, _ := ioutil.ReadAll(r.Body)
		data := appMetaData{}
		errFunc := yaml.Unmarshal(reqBody, &data)
		if errFunc != nil {
			logger.Warnf("Unable to unmarshall POST request: %v", errFunc)
			goto done
		}
		logger.Info("Unmarshalled POST request")

		// Validate metadata and add to metadata store
		errFunc = validateAppMetaData(data)
		if errFunc != nil {
			logger.Warnf("Failed to validate app metadata: %v", errFunc)
			goto done
		}
		select {
		case <-ctx.Done():
			errFunc = ctx.Err() // Context is done so don't perform more operations
			goto done
		default:
		}
		errFunc = dataStore.Add(data)
		if errFunc != nil {
			logger.Warnf("Failed to add app metadata: %v", errFunc)
			goto done
		}
		logger.Infof("Stored metadata for app '%v'", data.Title)
		logger.Infof("App metadata store contains %v entries", dataStore.TotalEntries())
	done:
		c <- errFunc
	}(ctx)
	select {
	case <-ctx.Done(): // End operation if ctx.Done() occurs
		return ctx.Err()
	case err := <-c:
		return err
	}
}

//
// Handler for creating new metadata entry
//
func createNewAppMetaData(w http.ResponseWriter, r *http.Request) {
	err := doCreateNewAppMetaData(r.Context(), r)
	if err != nil {
		logger.Warnf("createNewAppMetaData error: %v", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}
