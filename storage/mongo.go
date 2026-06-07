package storage

import (
	"context"
	"fmt"
	"time"

	"go-basics/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
	client *mongo.Client
}

func Connect(ctx context.Context, uri string) (*Store, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("ping mongodb: %w", err)
	}

	return &Store{client: client}, nil
}

func (s *Store) Close(ctx context.Context) error {
	return s.client.Disconnect(ctx)
}

func (s *Store) GetPlayers(ctx context.Context) ([]models.Player, error) {
	collection := s.client.Database("go-basics").Collection("player-data")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var players []models.Player
	for cursor.Next(ctx) {
		var player models.Player
		if err := cursor.Decode(&player); err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func (s *Store) GetPlayer(ctx context.Context, playerId string) (models.Player, error) {
	collection := s.client.Database("go-basics").Collection("player-data")
	var player models.Player
	err := collection.FindOne(ctx, bson.M{"player_id": playerId}).Decode(&player)
	if err != nil {
		return models.Player{}, err
	}
	return player, nil
}

func (s *Store) StorePlayer(ctx context.Context, player models.Player) error {
	collection := s.client.Database("go-basics").Collection("player-data")
	filter := bson.M{"player_id": player.PlayerId}
	update := bson.M{
		"$set": bson.M{
			"player_name": player.PlayerName,
			"player_team": player.PlayerTeam,
			"at_bats":     player.AtBats,
			"walks":       player.Walks,
			"strikeouts":  player.Strikeouts,
			"home_runs":   player.HomeRuns,
			"single":      player.Single,
			"double":      player.Double,
			"triple":      player.Triple,
			"home_run":    player.HomeRun,
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}
	_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}
	fmt.Println("Player stored or updated successfully")
	return nil
}

func (s *Store) DeletePlayer(ctx context.Context, playerId string) error {
	collection := s.client.Database("go-basics").Collection("player-data")
	_, err := collection.DeleteOne(ctx, bson.M{"player_id": playerId})
	if err != nil {
		return fmt.Errorf("delete player: %w", err)
	}
	fmt.Println("Player deleted successfully")
	return nil	
}