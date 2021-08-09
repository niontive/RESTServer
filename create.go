package main

import (
	"context"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

func dataStoreAdd(ctx context.Context, ch chan error, data appMetaData) {
	logger.Info("Adding to datastore")
	if ctx.Err() != nil {
		ch <- ctx.Err()
		return
	}
	err := dataStore.Add(data)
	logger.Info("Finished adding to datastore")
	ch <- err
}

//
// Add metadata entry to datastore
//
func doAddAppMetaData(ctx context.Context, data appMetaData) error {
	ch := make(chan error)

	go dataStoreAdd(ctx, ch, data)

	select {
	case <-ctx.Done():
		logger.Warn("Request context is closed")
		return ctx.Err()
	case err := <-ch:
		return err
	}
}

//
// Create and add new metadata entry via user request
//
func doCreateNewAppMetaData(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var err error
	logger.Info("Endpoint Hit: createNewAppMetaData")

	// Read and validate request into metadata object
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := appMetaData{}
	err = yaml.Unmarshal(reqBody, &data)
	if err != nil {
		logger.Warnf("Unable to unmarshall POST request: %v", err)
		goto done
	}
	logger.Info("Unmarshalled POST request")

	err = validateAppMetaData(data)
	if err != nil {
		logger.Warnf("Failed to validate app metadata: %v", err)
		goto done
	}

	// Add entry to datastore
	err = doAddAppMetaData(ctx, data)
	if err != nil {
		logger.Warnf("Failed to add app metadata: %v", err)
		goto done
	}
	logger.Infof("Stored metadata for app '%v'", data.Title)
	logger.Infof("App metadata store contains %v entries", dataStore.TotalEntries())

done:
	if err != nil {
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
	return
}

//
// Handler for creating new metadata entry
//
func createNewAppMetaData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	doCreateNewAppMetaData(ctx, w, r)
}
