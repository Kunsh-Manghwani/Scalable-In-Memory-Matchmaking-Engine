package internal

import (
    "testing"
)

func TestComputeGeohash(t *testing.T) {
    lat := 13.7563
    lon := 100.5018
    gh := computeGeohash(lat, lon)
    if len(gh) != 5 {
        t.Errorf("Expected geohash of length 5, got %d", len(gh))
    }
}

func TestCalculateScore(t *testing.T) {
    p1 := &Profile{ID: "1", Age: 28, Interests: []string{"music", "art"}}
    p2 := &Profile{ID: "2", Age: 30, Interests: []string{"music", "travel"}}
    score := calculateScore(p1, p2)
    if score <= 0 || score > 1 {
        t.Errorf("Expected score between 0 and 1, got %f", score)
    }
}

func TestSharedInterestScore(t *testing.T) {
    a := []string{"music", "art"}
    b := []string{"music", "travel"}
    s := sharedInterestScore(a, b)
    if s <= 0 {
        t.Errorf("Expected positive shared interest score, got %f", s)
    }
}

func TestProfileInsertionAndRetrieval(t *testing.T) {
    p := &Profile{ID: "testuser", Age: 25, Location: Location{Lat: 13.7, Lon: 100.5}, Interests: []string{"music"}}
    Profiles.Store(p.ID, p)

    val, ok := Profiles.Load(p.ID)
    if !ok {
        t.Errorf("Failed to retrieve profile by ID")
    }
    got := val.(*Profile)
    if got.ID != p.ID {
        t.Errorf("Retrieved wrong profile, expected %s, got %s", p.ID, got.ID)
    }
}

func TestExclusionFiltering(t *testing.T) {
    exclusions := &Exclusions{
        Matched:  map[string]bool{"u2": true},
        Blocked:  map[string]bool{"u3": true},
        Disliked: map[string]bool{},
    }
    ExclusionStore.Store("u1", exclusions)

    raw, _ := ExclusionStore.Load("u1")
    ex := raw.(*Exclusions)

    if !ex.Matched["u2"] || !ex.Blocked["u3"] {
        t.Errorf("Exclusions not set properly")
    }
}
