package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func registerHandler() *mux.Router  {
	router := mux.NewRouter()
	router.HandleFunc("/user",CreateUser).Methods("POST")
	router.HandleFunc("/user/{username}",Login).Methods("POST")

	return router
}


func main(){
  router := registerHandler()
  log.Fatal(http.ListenAndServe(":8080",router))
}