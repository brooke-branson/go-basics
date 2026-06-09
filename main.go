package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"go-basics/config"
	"go-basics/models"
	"go-basics/storage"
	"go-basics/transforms"
	"go-basics/validation"
	"log"
	"net/http"
	"time"
)

type app struct {
	store *storage.Store
}

func getHealthz(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	log.Println("Health check successful")
	json.NewEncoder(w).Encode(response)
}

func (app *app) getPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := app.store.GetPlayers(r.Context())
	if err != nil {
		log.Println("Error getting players:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(players)
}

func (app *app) getPlayer(w http.ResponseWriter, r *http.Request) {
	playerId := r.URL.Query().Get("player_id")
	player, err := app.store.GetPlayer(r.Context(), playerId)
	if err != nil {
		log.Println("Error getting player:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(player)
}

func (app *app) getTransformedPlayerXML(w http.ResponseWriter, r *http.Request) {
	playerId := r.URL.Query().Get("player_id")
	xmlOutput, err := app.store.GetTransformedPlayer(r.Context(), playerId)
	if err != nil {
		log.Println("Error getting transformed player:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xmlOutput))
}

func (app *app) addPlayer(w http.ResponseWriter, r *http.Request) {
	var player models.Player

	if r.Method != "POST" {
		log.Println("Only use POST method")
		http.Error(w, "Only use POST method", http.StatusMethodNotAllowed)
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&player); err != nil {
		log.Println("Error decoding player:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.ValidatePlayer(player); err != nil {
		log.Println("Error validating player:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := app.store.StorePlayer(r.Context(), player); err != nil {
		log.Println("Error storing player:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(player); err != nil {
		log.Println("Error encoding player:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (app *app) deletePlayer(w http.ResponseWriter, r *http.Request) {
	playerId := r.URL.Query().Get("player_id")
	if err := app.store.DeletePlayer(r.Context(), playerId); err != nil {
		log.Println("Error deleting player:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println("Player deleted successfully")
	json.NewEncoder(w).Encode(map[string]string{"message": "Player deleted successfully"})
}

func (app *app) transformPlayer(w http.ResponseWriter, r *http.Request) {
	// getTransformedPlayerXML takes a player ID and returns the transformed player in XML

	playerId := r.URL.Query().Get("player_id")
	player, err := app.store.GetPlayer(r.Context(), playerId)
	if err != nil {
		log.Println("Error getting player:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transformedPlayer := transforms.TransformPlayer(player)

	output, err := xml.MarshalIndent(transformedPlayer, "", "  ")

	if err != nil {
		log.Println("Error marshalling transformed player:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := app.store.StoreTransformedPlayer(r.Context(), playerId, string(output)); err != nil {
		log.Println("Error storing transformed player:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Transformed player stored successfully")
	w.Header().Set("Content-Type", "application/xml")
	w.Write(output)
}

func main() {
	cfg := config.LoadConfig()

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
	http.HandleFunc("/get-all-players", server.getPlayers)
	http.HandleFunc("/get-player", server.getPlayer)
	http.HandleFunc("/players/add", server.addPlayer)
	http.HandleFunc("/players/delete", server.deletePlayer)
	http.HandleFunc("/transform-player", server.transformPlayer)
	http.HandleFunc("/get-transformed-player", server.getTransformedPlayerXML)

	http.ListenAndServe(":"+cfg.Port, nil)
}
