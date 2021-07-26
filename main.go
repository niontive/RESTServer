package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const port = "10000"

var (
	dataStore = appMetaDataStore{store: make([]appMetaData, 0), dupTracker: make(map[string]bool)}
	logger    = logrus.New()
)

func createNewAppMetaData(w http.ResponseWriter, r *http.Request) {
	logger.Info("Endpoint Hit: createNewAppMetaData")
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := appMetaData{}
	err := yaml.Unmarshal(reqBody, &data)
	if err != nil {
		logger.Warnf("Unable to unmarhsall POST request: %v", err)
		return
	}
	logger.Info("Unmarshalled POST request")
	err = validateAppMetaData(data)
	if err != nil {
		logger.Warnf("Failed to validate app metadata: %v", err)
		return
	}
	err = dataStore.Add(data)
	if err != nil {
		logger.Warnf("Failed to add app metadata: %v", err)
		return
	}
	logger.Infof("Stored metadata for app '%v'", data.Title)
	logger.Infof("App metadata store contains %v entries", dataStore.TotalEntries())
}

func getAppMetaData(w http.ResponseWriter, r *http.Request) {
	logger.Info("EndpointHit: getAppMetaData")
	query := r.URL.Query()
	tmpMetaData := []appMetaData{}
	encoder := yaml.NewEncoder(w)
	for k, v := range query {
		var err error
		// Underscore character is treated as a space for certain fields
		v = replaceUnderscore(k, v)
		logger.Infof("Query key: %v. Query value: %v.", k, v)
		tmpMetaData, err = dataStore.Search(k, v)
		if err == nil {
			logger.Infof("Found metadata for key %v", k)
			for _, data := range tmpMetaData {
				encoder.Encode(data)
			}

		} else {
			logger.Warnf("Unable to find metadata for key '%v': %v", k, err)
		}
	}
	encoder.Close()
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
	router.HandleFunc("/createmetadata", createNewAppMetaData).Methods("POST")
	router.HandleFunc("/getmetadata", getAppMetaData).Methods("GET")
	logger.Fatal(http.ListenAndServe(":"+port, limit(router)))
}

func main() {
	logger.Info("Start REST server")
	handleRequests()
}
