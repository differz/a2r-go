package main

import (
	"fmt"
	"log"
	"net/http"

	"../../config"
	"../../data/phonebook"
	"../../manager"

	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/lib/pq"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server request, URL %s", r.URL.Path[1:])

	cfg := config.New()
	err := cfg.Init()
	if err != nil {
		log.Fatal(err)
	}
	phonebook.Get()
	manager := manager.New()
	store, err := manager.Create()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	_ = store
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
