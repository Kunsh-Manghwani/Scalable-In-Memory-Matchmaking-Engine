package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "matchmaking/internal"
)

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/profiles", internal.CreateProfile).Methods("POST")
    r.HandleFunc("/match/{id}", internal.GetMatches).Methods("GET")
    r.HandleFunc("/seed", internal.SeedProfiles).Methods("GET")
    r.HandleFunc("/profiles", internal.GetAllProfiles).Methods("GET")
    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
