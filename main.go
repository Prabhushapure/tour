package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Destination represents a travel destination
type Destination struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var db *sql.DB

func main() {
	initDB()

	router := mux.NewRouter()
	router.HandleFunc("/api/destinations", getDestinations).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server is running on port %s...\n", port)
	http.ListenAndServe(":"+port, router)
}

func initDB() {
	var err error
	db, err = sql.Open("mysql", "username:password@tcp(localhost:3306)/travelapp")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to the database")
}

func getDestinations(w http.ResponseWriter, r *http.Request) {
	destinations := fetchDestinationsFromDB()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(destinations)
}

func fetchDestinationsFromDB() []Destination {
	rows, err := db.Query("SELECT * FROM destinations")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var destinations []Destination
	for rows.Next() {
		var destination Destination
		err := rows.Scan(&destination.ID, &destination.Name, &destination.Description)
		if err != nil {
			log.Fatal(err)
		}
		destinations = append(destinations, destination)
	}

	return destinations
}
