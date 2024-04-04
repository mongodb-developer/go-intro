package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Scope struct {
	Project string
	Area    string
}

type Note struct {
	Title string
	Tags  []string
	Text  string
	Scope Scope
}

func main() {
	const serverAddr string = "127.0.0.1:8081"

	fmt.Println("Hola Caracola")

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HTTP Caracola"))
	})
	http.HandleFunc("POST /notes", createNote)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "Note: %+v", note)
}
