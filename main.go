package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const port = ":10000"

var (
	dataStore = appMetaDataStore{store: make([]appMetaData, 0)}
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
	dataStore.Add(data)
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
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/createmetadata", createNewAppMetaData).Methods("POST")
	myRouter.HandleFunc("/getmetadata", getAppMetaData).Methods("GET")
	logger.Fatal(http.ListenAndServe(port, myRouter))
}

func main() {
	logger.Info("Start REST server")
	handleRequests()
}
