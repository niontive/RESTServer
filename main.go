package main

import (
	"fmt"
	"log"
	"net/http"
)

const port = "10000"

type AppMetaData struct {
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Maintainers struct {
		Name  string `yaml:"name"`
		Email string `yaml:"email"`
	}
	Company     string `yaml:"company"`
	Website     string `yaml:"website"`
	Source      string `yaml:"source"`
	License     string `yaml:"license"`
	Description string `yaml:"description"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func main() {
	fmt.Println("Start server")
	handleRequests()
}
