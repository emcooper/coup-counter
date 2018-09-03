package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Game struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func GetGames(w http.ResponseWriter, r *http.Request) {
	log.Println("Getting games")
	games, err := store.GetGames()
	if err != nil {

	}
	json.NewEncoder(w).Encode(games)
}

func CreateGame(w http.ResponseWriter, r *http.Request) {}

func main() {
	log.Println("Welcome Coup-Counter")
  godotenv.Load()
  
	log.Printf("Opening db connection with database url: %v", os.Getenv("DATABASE_URL"))
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
  
	log.Printf("Listening on port %v", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}
