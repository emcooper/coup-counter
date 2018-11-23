package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Game struct {
	ID   int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Round struct {
	ID     int    `json:"id,omitempty"`
	Date   string `json:"date,omitempty"`
	GameId int    `json:"game_id,omitempty"`
}

type Player struct {
	ID   int `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CoupResult struct {
	ID             int `json:"id,omitempty"`
	RoundId        int    `json:"round_id,omitempty"`
	PlayerId       int    `json:"player_id,omitempty"`
	Winner         bool   `json:"winner,omitempty"`
	WinningCardOne string `json:"winning_card_one,omitempty"`
	WinningCardTwo string `json:"winning_card_two,omitempty"`
}

type CoupResultsWrapper struct {
	Results []CoupResult `json:"results,omitempty"`
}

func GetGames(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("Getting games")
	games, err := store.GetGames()
	if err != nil {

	}
	json.NewEncoder(w).Encode(games)
}

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("Getting players")
	players, err := store.GetPlayers()
	if err != nil {

	}
	json.NewEncoder(w).Encode(players)
}

func Migrate(w http.ResponseWriter, r *http.Request) {
	log.Println("Running migrations")

	err := store.Migrate()
	if err != nil {
		log.Println(err)
	}
}

func createCoupRound(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	log.Println("Adding Coup Round")

	var round Round
	round.Date = time.Now().Local().Format("2006-01-02")
	round.GameId = 1
	round, err := store.CreateRound(&round)
	if err != nil {
		log.Println(err)
	}

	// params := mux.Vars(r)
	results := CoupResultsWrapper{}
	err = json.NewDecoder(r.Body).Decode(&results)
	if err != nil {
		log.Println(err)
	}

	for _, result := range results.Results {
		result.RoundId = round.ID
		err = store.CreateCoupResult(&result)
	}
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w)
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
	router.HandleFunc("/migrate", Migrate).Methods("GET")
	router.HandleFunc("/games/coup", createCoupRound).Methods("POST")
	router.HandleFunc("/players", GetPlayers).Methods("GET")

	log.Printf("Listening on port %v", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), router))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}