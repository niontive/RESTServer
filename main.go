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

func doCreateNewAppMetaData(ctx context.Context, r *http.Request) error {
	c := make(chan error)
	go func(ctx context.Context) {
		logger.Info("Endpoint Hit: createNewAppMetaData")
		reqBody, _ := ioutil.ReadAll(r.Body)
		data := appMetaData{}
		errFunc := yaml.Unmarshal(reqBody, &data)
		if errFunc != nil {
			logger.Warnf("Unable to unmarshall POST request: %v", errFunc)
			goto done
		}
		logger.Info("Unmarshalled POST request")
		errFunc = validateAppMetaData(data)
		if errFunc != nil {
			logger.Warnf("Failed to validate app metadata: %v", errFunc)
			goto done
		}
		select {
		case <-ctx.Done():
			errFunc = ctx.Err()
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
	case <-ctx.Done():
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func createNewAppMetaData(w http.ResponseWriter, r *http.Request) {
	err := doCreateNewAppMetaData(r.Context(), r)
	if err != nil {
		logger.Warnf("createNewAppMetaData error: %v", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}

func doGetAppMetaData(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	c := make(chan bool)
	go func(ctx context.Context) {
		logger.Info("EndpointHit: getAppMetaData")
		var metaDataFound bool = false
		query := r.URL.Query()
		tmpMetaData := []appMetaData{}
		encoder := yaml.NewEncoder(w)
		for k, v := range query {
			var errFunc error
			select {
			case <-ctx.Done():
				return
			default:
			}
			// Underscore character is treated as a space for certain fields
			v = replaceUnderscore(k, v)
			logger.Infof("Query key: %v. Query value: %v.", k, v)
			tmpMetaData, errFunc = dataStore.Search(k, v)
			if errFunc == nil {
				logger.Infof("Found metadata for key %v", k)
				for _, data := range tmpMetaData {
					encoder.Encode(data)
					metaDataFound = true
				}

			} else {
				logger.Warnf("Unable to find metadata for key '%v': %v", k, errFunc)
			}
		}
		encoder.Close()
		c <- metaDataFound
	}(ctx)
	select {
	case <-ctx.Done():
		return ctx.Err()
	case result := <-c:
		if result == true {
			return nil
		} else {
			return errors.New("No metadata found")
		}
	}
}

func getAppMetaData(w http.ResponseWriter, r *http.Request) {
	err := doGetAppMetaData(r.Context(), w, r)
	if err != nil {
		logger.Warnf("getAppMetaData error: %v", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
	}
}

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
