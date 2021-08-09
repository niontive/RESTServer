package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const port = "10000"
const duration = 5 * time.Second // Timeout duration for APIs

var (
	dataStore = appMetaDataStore{store: make([]appMetaData, 0), dupTracker: make(map[string]bool)}
	logger    = logrus.New()
)

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
// Create endpoints and start server
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
