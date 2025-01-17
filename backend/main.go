package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Name struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	var err error
	db, err = sql.Open("postgres", "user=testusr password=testing dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Create table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS names (id SERIAL PRIMARY KEY, name TEXT NOT NULL)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("names table created successfully")

	router := mux.NewRouter()
	router.HandleFunc("/names", createName).Methods("POST")
	router.HandleFunc("/names", getNames).Methods("GET")
	router.HandleFunc("/names/{id}", updateName).Methods("PUT")
	router.HandleFunc("/names/{id}", deleteName).Methods("DELETE")

	// Add the CORS middleware
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"})

	fmt.Println("Listening and serving HTTP on :8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func createName(w http.ResponseWriter, r *http.Request) {
	var name Name
	err := json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var id int
	err = db.QueryRow(`INSERT INTO names (name) VALUES ($1) RETURNING id`, name.Name).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name.ID = id
	json.NewEncoder(w).Encode(name)
}

func getNames(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, name FROM names`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var names []Name
	for rows.Next() {
		var name Name
		err := rows.Scan(&name.ID, &name.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		names = append(names, name)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(names)
}

func updateName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var name Name
	err := json.NewDecoder(r.Body).Decode(&name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec(`UPDATE names SET name=$1 WHERE id=$2`, name.Name, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(name)
}

func deleteName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("Received DELETE request for id: %s\n", id) // Debug statement

	result, err := db.Exec(`DELETE FROM names WHERE id=$1`, id)
	if err != nil {
		fmt.Printf("Error executing DELETE: %s\n", err) // Debug statement
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error fetching rows affected: %s\n", err) // Debug statement
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		fmt.Printf("No rows were deleted for id: %s\n", id) // Debug statement
		http.Error(w, "No rows were deleted", http.StatusNotFound)
		return
	}

	fmt.Printf("Successfully deleted id: %s\n", id) // Debug statement
	w.WriteHeader(http.StatusNoContent)
}
