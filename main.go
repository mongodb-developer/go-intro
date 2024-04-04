package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

var mdbClient *mongo.Client

func main() {
	const serverAddr string = "127.0.0.1:8081"
	// TODO: Replace with your connection string
	const connStr string = "mongodb+srv://yourusername:yourpassword@notekeeper.xxxxxx.mongodb.net/?retryWrites=true&w=majority&appName=NoteKeeper"
	done := make(chan struct{})

	fmt.Println("Hola Caracola")

	ctxBg := context.Background()
	var err error
	mdbClient, err = mongo.Connect(ctxBg, options.Client().ApplyURI(connStr))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = mdbClient.Disconnect(ctxBg); err != nil {
			panic(err)
		}
	}()

	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HTTP Caracola"))
	})
	router.HandleFunc("POST /notes", createNote)

	server := http.Server{
		Addr:    serverAddr,
		Handler: router,
	}
	server.RegisterOnShutdown(func() {
		defer func() {
			done <- struct{}{}
		}()
		fmt.Println("Signal shutdown")
		time.Sleep(5 * time.Second)
	})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("Server shutdown error: %v", err)
		}
	}()
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Panicf("HTTP server error %v\n", err)
	}
	<-done
}

func createNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&note); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	notesCollection := mdbClient.Database("NoteKeeper").Collection("Notes")
	result, err := notesCollection.InsertOne(r.Context(), note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Note: %+v", note)
	log.Printf("Id: %v", result.InsertedID)
}
