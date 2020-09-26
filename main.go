package main

import (
	"fmt"
	"golang-api/controllers"
	"golang-api/kafka"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/getbyid/{id}+{limit}", controllers.ReturnData)
	myRouter.HandleFunc("/addData", controllers.AddNewData).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Println("Rest API using Mux Routers")
	kafka.Consume()
	handleRequests()
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
