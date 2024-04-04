package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const serverAddr string = "127.0.0.1:8081"

	fmt.Println("Hola Caracola")

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HTTP Caracola"))
	})
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}
