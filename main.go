package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const port = "localhost:10000"
const duration = 5 * time.Second     // Timeout duration for APIs
const shutdownTime = 1 * time.Second // Timeout duration for shutdown

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

	srv := &http.Server{
		Addr:         port,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      limit(router),
	}
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Fatalf("Server stopped: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err != nil {
		logger.Warnf("Shutdown error: %v", err)
	}
	logger.Infof("Shutting Down")
	os.Exit(0)
}

func main() {
	logger.Info("Start REST server")
	handleRequests()
}
