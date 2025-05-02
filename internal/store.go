package internal

import (
    "sync"
)

var (
    Profiles      = sync.Map{} // id → *Profile
    GeohashIndex = sync.Map{} // geohash → map[id]*Profile
    MatchScores  = sync.Map{} // id → []MatchScore
    ExclusionStore = sync.Map{} // id → *Exclusions
)
