package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

const jsonDir = "./data/"

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/set/{endpoint}", setEndpointData).Methods("POST")
	router.HandleFunc("/data/{endpoint}", getEndpointData).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func setEndpointData(w http.ResponseWriter, r *http.Request) {
	endpoint := mux.Vars(r)["endpoint"]
	fileName := fmt.Sprintf("%s.json", endpoint)
	filePath := path.Join(jsonDir, fileName)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read request body"))
		return
	}

	err = ioutil.WriteFile(filePath, body, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to create JSON file"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("JSON file created"))
}

func getEndpointData(w http.ResponseWriter, r *http.Request) {
	endpoint := mux.Vars(r)["endpoint"]
	fileName := fmt.Sprintf("%s.json", endpoint)
	filePath := path.Join(jsonDir, fileName)

	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to read file"))
		return
	}

	var jsonData interface{}
	err = json.Unmarshal(file, &jsonData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to parse JSON"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonData)
}
