package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"context"
	"net/http"
	"go-basics/models"
	"go-basics/storage"
	"go-basics/transforms"
	"go-basics/config"
	"go-basics/validation"
	"time"
)

type app struct {
	store *storage.Store
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}
	json.NewEncoder(w).Encode(response)
}

func (app *app) getPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := app.store.GetPlayers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

func (app *app) getPlayer(w http.ResponseWriter, r *http.Request){
	playerId := r.URL.Query().Get("player_id")
	player, err := app.store.GetPlayer(r.Context(), playerId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(player)
}

func (app *app) addPlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player

	if r.Method != "POST" {
		http.Error(w, "Only use POST method", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.ValidatePlayer(player); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := app.store.StorePlayer(r.Context(), player); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(player)
}

func (app *app) deletePlayer(w http.ResponseWriter, r *http.Request) {
	playerId := r.URL.Query().Get("player_id")
	if err := app.store.DeletePlayer(r.Context(), playerId); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Player deleted successfully"})
}

func (app *app) getPlayerTransformed(w http.ResponseWriter, r *http.Request) {
	// getPlayersTransformed takes a pre-formatted ojbject in JSON
	// and converts it to XML
	var player models.Player

	if r.Method != "POST" {
		http.Error(w, "Only use POST method", http.StatusMethodNotAllowed)
		return
	}

	json.NewDecoder(r.Body).Decode(&player)
	transformedPlayer := transforms.TransformPlayer(player)

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

func (app *app) getTransformedPlayerXML(w http.ResponseWriter, r *http.Request) {
	// getTransformedPlayerXML takes a player ID and returns the transformed player in XML

	playerId := r.URL.Query().Get("player_id")
	player, err := app.store.GetPlayer(r.Context(), playerId)
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	transformedPlayer := transforms.TransformPlayer(player)

	output, err := xml.MarshalIndent(transformedPlayer, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := app.store.StoreTransformedPlayer(r.Context(), playerId, string(output)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(output)
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	store, err := storage.Connect(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close(context.Background())

	fmt.Println("Connected to MongoDB")

	server := &app{store: store}

	fmt.Println("Listening on port ", cfg.Port)

	http.HandleFunc("/healthz", getHealthz)
	http.HandleFunc("/players", server.getPlayers)
	http.HandleFunc("/players/lookup", server.getPlayer)
	http.HandleFunc("/players/add", server.addPlayer)
	http.HandleFunc("/players/delete", server.deletePlayer)
	http.HandleFunc("/transformed-players/json", server.getPlayerTransformed)
	http.HandleFunc("/transformed-player/xml", server.getTransformedPlayerXML)
	http.HandleFunc("/transformed-players/xml", server.getPlayerTransformed)

	http.ListenAndServe(":8080", nil)
}