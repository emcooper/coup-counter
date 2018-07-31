package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)


type Game struct {
    ID        string   `json:"id,omitempty"`
    Name string   `json:"name,omitempty"`
}

var games []Game




func GetPeople(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(games)
}
func GetPerson(w http.ResponseWriter, r *http.Request) {}
func CreatePerson(w http.ResponseWriter, r *http.Request) {}
func DeletePerson(w http.ResponseWriter, r *http.Request) {}

// our main function
func main() {
    router := mux.NewRouter()
    games = append(games, Game{ID: "1", Name: "Coup"})
    games = append(games, Game{ID: "2", Name: "7 Wonders"})
    router.HandleFunc("/people", GetPeople).Methods("GET")
    router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8000", router))
}
