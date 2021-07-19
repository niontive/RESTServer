package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const port = ":10000"

func createNewAppMetaData(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
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
