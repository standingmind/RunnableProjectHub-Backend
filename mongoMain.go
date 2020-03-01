package main

import (
	"context"
	"encoding/json"
	"fmt"
	"helper"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/getProjects", getProjects).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

