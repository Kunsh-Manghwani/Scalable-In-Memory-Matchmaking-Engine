package internal

type Location struct {
    Lat float64 `json:"lat"`
    Lon float64 `json:"lon"`
}

type Profile struct {
    ID         string   `json:"id"`
    Age        int      `json:"age"`
    Gender     string   `json:"gender"`
    Location   Location `json:"location"`
    Interests  []string `json:"interests"`
    LookingFor string   `json:"looking_for,omitempty"`
}

type MatchScore struct {
    ProfileID string  `json:"profile_id"`
    Score     float64 `json:"score"`
}

type Exclusions struct {
    Matched  map[string]bool
    Blocked  map[string]bool
    Disliked map[string]bool
}
