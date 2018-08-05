package main

import (
    "encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "database/sql"
    _"github.com/lib/pq"
    // "io/ioutil"
)


type Game struct {
    ID        string   `json:"id,omitempty"`
    Name string   `json:"name,omitempty"`
}






func GetPeople(w http.ResponseWriter, r *http.Request) {
  games, err := store.GetGames()
  if err != nil {

  }
  json.NewEncoder(w).Encode(games)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {}
func CreatePerson(w http.ResponseWriter, r *http.Request) {}
func DeletePerson(w http.ResponseWriter, r *http.Request) {}

// our main function
func main() {
    connString := "dbname=coup_counter_development sslmode=disable"
    db, err := sql.Open("postgres", connString)

    if err != nil {
      panic(err)
    }
    err = db.Ping()

    if err != nil {
      panic(err)
    }

    InitStore(&dbStore{db: db})
    router := mux.NewRouter()

    router.HandleFunc("/people", GetPeople).Methods("GET")
    router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
    router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
    router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
    log.Fatal(http.ListenAndServe(":8000", router))
}
