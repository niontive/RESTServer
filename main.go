package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const port = "10000"
const duration = 5 * time.Second // Timeout duration for APIs

var (
	dataStore = appMetaDataStore{store: make([]appMetaData, 0), dupTracker: make(map[string]bool)}
	logger    = logrus.New()
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

//
// Homepage that displays stored metadata
//
func homePage(w http.ResponseWriter, r *http.Request) {
	logger.Info("Endpoint Hit: homePage")
	if dataStore.TotalEntries() == 0 {
		fmt.Fprintf(w, "No application meta data stored!")
	} else {
		// Obtain and display application titles
		fmt.Fprintf(w, "Available app meta data:\n")
		titles := dataStore.GetAppTitles()
		for _, element := range titles {
			fmt.Fprintf(w, element+"\n")
		}
	}
}

//
// Create routes for endpoints
//
func handleRequests() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homePage)
	router.Handle("/createmetadata", http.TimeoutHandler(http.HandlerFunc(createNewAppMetaData),
		duration, "Timeout createmetadata\n")).Methods("POST")
	router.Handle("/getmetadata", http.TimeoutHandler(http.HandlerFunc(getAppMetaData),
		duration, "Timeout getmetadata\n")).Methods("GET")
	logger.Fatal(http.ListenAndServe(":"+port, limit(router)))
}

func main() {
	logger.Info("Start REST server")
	handleRequests()
}
