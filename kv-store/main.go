package main

import (
	"fmt"
	"kv-store/api"
	"kv-store/store"
	"log"
	"net/http"
)

func main() {
	s := store.New()
	h := api.NewHandler(s)

	http.HandleFunc("/keys", h.HandleKeys)
	http.HandleFunc("/keys/", h.HandleKey)

	addr := ":8080"
	fmt.Printf("kv-store server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
