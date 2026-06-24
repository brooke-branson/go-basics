package mlbstats

import (
	"context"
	"testing"
)

// Hits the real MLB API — needs network access.
func TestFetchPlayer_Ohtani(t *testing.T) {
	client := NewClient()

	player, err := client.FetchPlayer(context.Background(), "660271", 2024)
	if err != nil {
		t.Fatalf("FetchPlayer failed: %v", err)
	}

	if player.PlayerName != "Shohei Ohtani" {
		t.Errorf("expected Shohei Ohtani, got %q", player.PlayerName)
	}

	t.Logf("player: %+v", player)
}

func TestFetchPlayer_Tatis(t *testing.T) {
	client := NewClient()

	player, err := client.FetchPlayer(context.Background(), "665487", 2026)
	if err != nil {
		t.Fatalf("FetchPlayer failed: %v", err)
	}

	if player.PlayerName != "Fernando Tatis Jr." {
		t.Errorf("expected Fernando Tatis Jr., got %q", player.PlayerName)
	}

	t.Logf("player: %+v", player)
}

func TestFetchPlayer_Machado(t *testing.T) {
	client := NewClient()

	player, err := client.FetchPlayer(context.Background(), "592518", 2026)
	if err != nil {
		t.Fatalf("FetchPlayer failed: %v", err)
	}

	if player.PlayerName != "Aaron Judge" {
		t.Errorf("expected Aaron Judge, got %q", player.PlayerName)
	}

	t.Logf("player: %+v", player)
}

func TestFetchPlayer_Judge(t *testing.T) {
	client := NewClient()

	player, err := client.FetchPlayer(context.Background(), "592450", 2026)
	if err != nil {
		t.Fatalf("FetchPlayer failed: %v", err)
	}

	if player.PlayerName != "Aaron Judge" {
		t.Errorf("expected Aaron Judge, got %q", player.PlayerName)
	}

	t.Logf("player: %+v", player)
}

func TestFetchPlayer_Harper(t *testing.T) {
	client := NewClient()

	player, err := client.FetchPlayer(context.Background(), "547180", 2026)
	if err != nil {
		t.Fatalf("FetchPlayer failed: %v", err)
	}

	if player.PlayerName != "Bryce Harper" {
		t.Errorf("expected Bryce Harper, got %q", player.PlayerName)
	}

	t.Logf("player: %+v", player)
}
