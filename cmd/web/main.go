package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/artist", artistByNumber)
	mux.HandleFunc("/filters", filters)
	mux.HandleFunc("/search", search)

	fileServer := http.FileServer(http.Dir("./templates/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Server go http://localhost:4000")
	err := http.ListenAndServe(":4000", mux)
	log.Println(err)
	return
}
