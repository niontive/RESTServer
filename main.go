package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const port = ":10000"

var logger = logrus.New()
var AppMetaDataStore = []AppMetaData{}

func createNewAppMetaData(w http.ResponseWriter, r *http.Request) {
	logger.Info("Endpoint Hit: createNewAppMetaData")
	reqBody, _ := ioutil.ReadAll(r.Body)
	data := AppMetaData{}
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
	AppMetaDataStore = append(AppMetaDataStore, data)
	logger.Infof("Stored metadata for app '%v'", data.Title)
	logger.Infof("App metadata store contains %v entries", len(AppMetaDataStore))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	logger.Info("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/appmetadata", createNewAppMetaData).Methods("POST")
	log.Fatal(http.ListenAndServe(port, myRouter))
}

func main() {
	logger.Info("Start REST server")
	handleRequests()
}
