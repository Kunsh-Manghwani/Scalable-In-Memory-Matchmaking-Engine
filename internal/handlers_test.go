package internal

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"sync"
    "github.com/gorilla/mux"
)

func setupRouter() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/profiles", CreateProfile).Methods("POST")
    r.HandleFunc("/profiles", GetAllProfiles).Methods("GET")
    r.HandleFunc("/match/{id}", GetMatches).Methods("GET")
    r.HandleFunc("/seed", SeedProfiles).Methods("GET")
    return r
}

func TestCreateAndGetProfile(t *testing.T) {
    router := setupRouter()
    payload := []byte(`{
        "id":"user123",
        "age":28,
        "gender":"F",
        "location":{"lat":13.7563,"lon":100.5018},
        "interests":["music","art"]
    }`)

    req, _ := http.NewRequest("POST", "/profiles", bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    if resp.Code != http.StatusCreated {
        t.Errorf("Expected status 201, got %d", resp.Code)
    }

    req, _ = http.NewRequest("GET", "/profiles", nil)
    resp = httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    if resp.Code != http.StatusOK {
        t.Errorf("Expected status 200, got %d", resp.Code)
    }
}

func TestSeedProfiles(t *testing.T) {
    router := setupRouter()
    req, _ := http.NewRequest("GET", "/seed?count=10", nil)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    if resp.Code != http.StatusOK {
        t.Errorf("Expected status 200 on seed, got %d", resp.Code)
    }
}

func TestGetMatches(t *testing.T) {
    router := setupRouter()

    // Create base user
    baseUser := Profile{
        ID:        "user999",
        Age:       25,
        Gender:    "F",
        Location:  Location{Lat: 13.7, Lon: 100.5},
        Interests: []string{"music", "travel"},
    }
    Profiles.Store(baseUser.ID, &baseUser)
    gh := computeGeohash(baseUser.Location.Lat, baseUser.Location.Lon)
    raw, _ := GeohashIndex.LoadOrStore(gh, &sync.Map{})
    quadrant := raw.(*sync.Map)
    quadrant.Store(baseUser.ID, &baseUser)
    precomputeMatches(&baseUser)

    // Request matches
    req, _ := http.NewRequest("GET", "/match/user999", nil)
    resp := httptest.NewRecorder()
    router.ServeHTTP(resp, req)

    if resp.Code != http.StatusOK {
        t.Errorf("Expected status 200 on match, got %d", resp.Code)
    }

    var matches []MatchScore
    if err := json.Unmarshal(resp.Body.Bytes(), &matches); err != nil {
        t.Errorf("Failed to parse match response: %v", err)
    }
}
