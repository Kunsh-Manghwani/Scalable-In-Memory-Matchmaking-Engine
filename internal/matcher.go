package internal

import (
    "math"
    "sort"
    "sync"
    "github.com/mmcloughlin/geohash"
)

func computeGeohash(lat, lon float64) string {
    return geohash.EncodeWithPrecision(lat, lon, 5)
}

func precomputeMatches(newProfile *Profile) {
    gh := computeGeohash(newProfile.Location.Lat, newProfile.Location.Lon)
    raw, _ := GeohashIndex.LoadOrStore(gh, &sync.Map{})
    quadrant := raw.(*sync.Map)

    var scores []MatchScore

    quadrant.Range(func(_, v interface{}) bool {
        other := v.(*Profile)
        if other.ID == newProfile.ID {
            return true
        }
        score := calculateScore(newProfile, other)
        scores = append(scores, MatchScore{ProfileID: other.ID, Score: score})

        // Update other's match list
        updateMatchList(other.ID, MatchScore{ProfileID: newProfile.ID, Score: score})
        return true
    })

    sort.Slice(scores, func(i, j int) bool { return scores[i].Score > scores[j].Score })
    if len(scores) > 5 {
        scores = scores[:5]
    }
    MatchScores.Store(newProfile.ID, scores)
}

func calculateScore(a, b *Profile) float64 {
    ageScore := 1 - math.Min(math.Abs(float64(a.Age-b.Age))/50, 1)
    interestScore := sharedInterestScore(a.Interests, b.Interests)
    return 0.4*ageScore + 0.6*interestScore
}

func sharedInterestScore(a, b []string) float64 {
    set := make(map[string]bool)
    for _, item := range a {
        set[item] = true
    }
    shared := 0
    for _, item := range b {
        if set[item] {
            shared++
        }
    }
    if len(a)+len(b) == 0 {
        return 0
    }
    return float64(shared) / float64(len(a))
}

func updateMatchList(userID string, match MatchScore) {
    raw, _ := MatchScores.LoadOrStore(userID, []MatchScore{})
    list := raw.([]MatchScore)
    list = append(list, match)
    sort.Slice(list, func(i, j int) bool { return list[i].Score > list[j].Score })
    if len(list) > 5 {
        list = list[:5]
    }
    MatchScores.Store(userID, list)
}
