package main

import (
	"database/sql"
)

type Store interface {
	CreateGame(game *Game) error
	CreateRound(round *Round) (Round, error)
	CreateCoupResult(result *CoupResult) error 
	GetGames() ([]*Game, error)
	Migrate() error
}

// The `dbStore` struct will implement the `Store` interface
// It also takes the sql DB connection object, which represents
// the database connection.
type dbStore struct {
	db *sql.DB
}

func (store *dbStore) CreateGame(game *Game) error {
	// 'Bird' is a simple struct which has "species" and "description" attributes
	// THe first underscore means that we don't care about what's returned from
	// this insert query. We just want to know if it was inserted correctly,
	// and the error will be populated if it wasn't
	_, err := store.db.Query("INSERT INTO games(name) VALUES ($1)", game.Name)
	return err
}

func (store *dbStore) GetGames() ([]*Game, error) {
	// Query the database for all birds, and return the result to the
	// `rows` object
	rows, err := store.db.Query("SELECT name from games")
	// We return incase of an error, and defer the closing of the row structure
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create the data structure that is returned from the function.
	// By default, this will be an empty array of birds
	games := []*Game{}
	for rows.Next() {
		// For each row returned by the table, create a pointer to a bird,
		game := &Game{}
		// Populate the `Species` and `Description` attributes of the bird,
		// and return incase of an error
		if err := rows.Scan(&game.Name); err != nil {
			return nil, err
		}
		// Finally, append the result to the returned array, and repeat for
		// the next row
		games = append(games, game)
	}
	return games, nil
}

func (store *dbStore) CreateRound(round *Round) (Round, error) {
	var insertId int
	err := store.db.QueryRow("INSERT INTO rounds(date, game_id) VALUES ($1, $2) RETURNING id", round.Date, round.GameId).Scan(&insertId)
	if err != nil {
		return *round, err
	}
	
	round.ID = insertId
	return *round, err
}

func (store *dbStore) CreateCoupResult(result *CoupResult) error {
	_, err := store.db.Query("INSERT INTO coup_results(round_id, player_id, winner, winning_card_one, winning_card_two) VALUES ($1, $2, $3, $4, $5)", result.RoundId, result.PlayerId, result.Winner, result.WinningCardOne, result.WinningCardTwo)
	return err
}

func (store *dbStore) Migrate() error {
	_, err := store.db.Query("CREATE TABLE IF NOT EXISTS games (id SERIAL PRIMARY KEY, name VARCHAR(255), created_at DATE NOT NULL DEFAULT CURRENT_DATE)")
	if err != nil {
		return err
	}
	_, err = store.db.Query("CREATE TABLE IF NOT EXISTS rounds (id SERIAL PRIMARY KEY, date DATE DEFAULT NULL, game_id INTEGER REFERENCES games(id), created_at DATE NOT NULL DEFAULT CURRENT_DATE)")
	if err != nil {
		return err
	}
	_, err = store.db.Query("CREATE TABLE IF NOT EXISTS players (id SERIAL PRIMARY KEY, name VARCHAR(255), created_at DATE NOT NULL DEFAULT CURRENT_DATE)")
	if err != nil {
		return err
	}
	_, err = store.db.Query("CREATE TABLE IF NOT EXISTS coup_results (id SERIAL PRIMARY KEY, round_id INTEGER REFERENCES rounds(id), player_id INTEGER REFERENCES players(id), winner BOOLEAN, winning_card_one VARCHAR(255), winning_card_two VARCHAR(255), created_at DATE NOT NULL DEFAULT CURRENT_DATE)")
	return err
}

// The store variable is a package level variable that will be available for
// use throughout our application code
var store Store

/*
We will need to call the InitStore method to initialize the store. This will
typically be done at the beginning of our application (in this case, when the server starts up)
This can also be used to set up the store as a mock, which we will be observing
later on
*/
func InitStore(s Store) {
	store = s
}
