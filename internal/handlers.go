package internal

import (
    "encoding/json"
    "math/rand"
    "net/http"
    "strconv"
    "sync"
    "github.com/gorilla/mux"
)

func CreateProfile(w http.ResponseWriter, r *http.Request) {
    var p Profile
    if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
        http.Error(w, "Invalid payload", http.StatusBadRequest)
        return
    }

    Profiles.Store(p.ID, &p)
    gh := computeGeohash(p.Location.Lat, p.Location.Lon)
    raw, _ := GeohashIndex.LoadOrStore(gh, &sync.Map{})
    quadrant := raw.(*sync.Map)
    quadrant.Store(p.ID, &p)

    precomputeMatches(&p)

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(p)
}

func GetMatches(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]
    raw, ok := MatchScores.Load(id)
    if !ok {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }
    matches := raw.([]MatchScore)

    rawEx, _ := ExclusionStore.LoadOrStore(id, &Exclusions{
        Matched:  make(map[string]bool),
        Blocked:  make(map[string]bool),
        Disliked: make(map[string]bool),
    })
    ex := rawEx.(*Exclusions)

    var result []MatchScore
    for _, m := range matches {
        if ex.Matched[m.ProfileID] || ex.Blocked[m.ProfileID] || ex.Disliked[m.ProfileID] {
            continue
        }
        result = append(result, m)
        if len(result) >= 5 {
            break
        }
    }

    json.NewEncoder(w).Encode(result)
}

func SeedProfiles(w http.ResponseWriter, r *http.Request) {
    countStr := r.URL.Query().Get("count")
    count, _ := strconv.Atoi(countStr)
    if count <= 0 {
        count = 100
    }
    interests := []string{"music", "art", "travel", "sports", "movies"}

    for i := 0; i < count; i++ {
        p := &Profile{
            ID:        "user" + strconv.Itoa(rand.Intn(100000)),
            Age:       rand.Intn(50) + 18,
            Gender:    []string{"M", "F"}[rand.Intn(2)],
            Location:  Location{Lat: 13.7 + rand.Float64(), Lon: 100.5 + rand.Float64()},
            Interests: []string{interests[rand.Intn(len(interests))]},
        }
        Profiles.Store(p.ID, p)
        gh := computeGeohash(p.Location.Lat, p.Location.Lon)
        raw, _ := GeohashIndex.LoadOrStore(gh, &sync.Map{})
        quadrant := raw.(*sync.Map)
        quadrant.Store(p.ID, p)

        precomputeMatches(p)
    }

    w.Write([]byte("Seeded profiles"))
}

func GetAllProfiles(w http.ResponseWriter, r *http.Request) {
    ageFilter := r.URL.Query().Get("age")
    geohashFilter := r.URL.Query().Get("geohash")

    var allProfiles []*Profile

    Profiles.Range(func(key, value interface{}) bool {
        p := value.(*Profile)

        // Age filter
        if ageFilter != "" {
            age, err := strconv.Atoi(ageFilter)
            if err != nil || p.Age != age {
                return true
            }
        }

        // Geohash filter
        if geohashFilter != "" {
            gh := computeGeohash(p.Location.Lat, p.Location.Lon)
            if gh != geohashFilter {
                return true
            }
        }

        allProfiles = append(allProfiles, p)
        return true
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(allProfiles)
}

