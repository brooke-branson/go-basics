package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"go-basics/models"
	"go-basics/transforms"
)

func getHealthz(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	response := models.Player{
		PlayerId: "1",
		PlayerName: "Fernando Tatis Jr.",
		PlayerTeam: "Padres",
	}
	json.NewEncoder(w).Encode(response)
}

func getTransformedPlayers(w http.ResponseWriter, r *http.Request) {

	var rawPlayer transforms.RawPlayer

	if r.Method != "POST" {
		http.Error(w, "Only use POST method", http.StatusMethodNotAllowed)
		return
	}

	json.NewDecoder(r.Body).Decode(&rawPlayer)
	transformedPlayer := transforms.TransformPlayer(rawPlayer)

	switch r.URL.Path {
	case "/transformed-players/json":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(transformedPlayer)

	case "/transformed-players/xml":
		w.Header().Set("Content-Type", "application/xml")
		output, err := xml.MarshalIndent(transformedPlayer, "", "  ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/xml")
		w.Write(output)
		
	default:
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
}
	
func main() {
	fmt.Println("Listening on port 8080")

	http.HandleFunc("/healthz", getHealthz)
	http.HandleFunc("/players", getPlayers)
	http.HandleFunc("/transformed-players/json", getTransformedPlayers)
	http.HandleFunc("/transformed-players/xml", getTransformedPlayers)

	http.ListenAndServe(":8080", nil)
}