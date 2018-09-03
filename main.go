package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	// "io/ioutil"
)

type Game struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func GetGames(w http.ResponseWriter, r *http.Request) {
	games, err := store.GetGames()
	if err != nil {

	}
	json.NewEncoder(w).Encode(games)
}

func CreateGame(w http.ResponseWriter, r *http.Request) {}

// our main function
func main() {
	log.Println("Welcome Coup-Counter")
	// connString := "dbname=coup_counter_development sslmode=disable"
	// db, err := sql.Open("postgres", connString)
	log.Println("Opening db connection")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Println("Panicing")
		panic(err)
	}
	log.Println("Pinging db")

	err = db.Ping()

	if err != nil {
		log.Println("Panicing")
		panic(err)
	}

	log.Println("Initializing store")

	InitStore(&dbStore{db: db})
	router := mux.NewRouter()

	router.HandleFunc("/games", GetGames).Methods("GET")
	router.HandleFunc("/", GetGames).Methods("GET")
	router.HandleFunc("/games", CreateGame).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
