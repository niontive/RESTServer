package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

const port = ":10000"

func createNewAppMetaData(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
	data := AppMetaData{}
	err := yaml.Unmarshal(reqBody, &data)
	if err != nil {
		log.Fatalf("Unable to unmarhsall POST request: %v", err)
	}
	fmt.Println("Unmarshalled POST request")
	err = validateAppMetaData(data)
	if err != nil {
		log.Fatalf("Failed to validate app meta data")
	}
	fmt.Println("Success!")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/appmetadata", createNewAppMetaData).Methods("POST")
	log.Fatal(http.ListenAndServe(port, myRouter))
}

func main() {
	fmt.Println("Start server")
	handleRequests()
}
