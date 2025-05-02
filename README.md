# Scalable In-Memory Matchmaking Engine

## 🚀 Project Overview

This project is a scalable, high-performance in-memory matchmaking engine for a dating app. It allows users to register profiles and retrieve their top 5 most compatible matches in real time. The system uses geohashing, a weighted scoring algorithm, and efficient in-memory data structures to deliver fast match lookups without needing a database.

---

## 🏗 Architecture Decisions

- **Language**: Go — chosen for its performance, concurrency support, and simple HTTP handling.
- **In-Memory Storage**: `sync.Map` is used for thread-safe storage of user profiles, geohash buckets, precomputed match scores, and exclusion lists.
- **Geohashing**: Profiles are bucketed using geohash (precision 5) to narrow down nearby match candidates efficiently.
- **API Framework**: Gorilla Mux router handles REST API endpoints.
- **Exclusion Handling**: Each user has tracked exclusions (matched, blocked, disliked) to avoid returning inappropriate matches.

---

## ⚙️ Pre-Computation Design

When a profile is registered:
1. The system calculates its geohash quadrant.
2. It finds all existing users in the same quadrant.
3. For each, it calculates a **match score** based on:
   - Age similarity  
   - Number of shared interests  
   - Location proximity (via geohash filtering)
4. It stores the top 5 matches for the new user and updates the match scores of existing users to include the newcomer.

This approach ensures that match retrieval is constant-time (`O(1)`), since results are already precomputed during registration.

---

## 💡 What I’d Do Differently in Production

- **Persistent Storage**: Replace in-memory maps with a database (PostgreSQL) or distributed cache (Redis).
- **Background Jobs**: Move pre-computation to background workers using a job queue.
- **Better Geospatial Search**: Use Redis GEO or PostGIS for more accurate and scalable location filtering.
- **Advanced Filtering**: Add gender preference, distance range, and interest weighting.
- **Scalability & Monitoring**: Use Kubernetes, health checks, metrics, horizontal scaling, and tracing.
- **API Security**: Add authentication, rate limiting, and payload validation.

---

## ⚙ Setup and Run

1. Install Go:
```bash
go version
````

2. Run the server:

```bash
go run cmd/main.go
```

---

## ✅ Run Tests

```bash
go test ./...
```

---

## 🌱 Seed Dummy Profiles

Populate 1000 fake profiles:

```bash
curl http://localhost:8080/seed?count=1000
```

---

## 🔍 Get All Profiles

```bash
curl http://localhost:8080/profiles
```

---

## 📬 Get Top 5 Matches

```bash
curl http://localhost:8080/match/user123
```

---

## 📚 API Documentation

### 📥 POST /profiles

Register a new user profile.

**Request Body (JSON)**:

```json
{
  "id": "user123",
  "age": 28,
  "gender": "F",
  "location": { "lat": 13.7563, "lon": 100.5018 },
  "interests": ["music", "art", "travel"]
}
```

**Response (201 Created)**:

```json
{
  "id": "user123",
  "age": 28,
  "gender": "F",
  "location": { "lat": 13.7563, "lon": 100.5018 },
  "interests": ["music", "art", "travel"]
}
```

---

### 📤 GET /profiles

Get all user profiles.

**Query Parameters (optional)**:

* `age` → filter by age → `/profiles?age=28`
* `geohash` → filter by geohash → `/profiles?geohash=xxxxx`

**Response**:

```json
[
  {
    "id": "user123",
    "age": 28,
    "gender": "F",
    "location": { "lat": 13.7563, "lon": 100.5018 },
    "interests": ["music", "art", "travel"]
  },
  ...
]
```

---

### 💘 GET /match/{id}

Get top 5 matches for a user.

**Example**:

```
GET /match/user123
```

**Response**:

```json
[
  { "profile_id": "user456", "score": 0.85 },
  { "profile_id": "user789", "score": 0.80 }
]
```

---

### 🌱 GET /seed

Bulk insert dummy profiles.

**Example**:

```
GET /seed?count=1000
```

**Response**:

```
Seeded profiles
```

---

## 🧪 Example Testing Commands

1. **Run tests:**

```bash
go test ./...
```

2. **Manually create a profile:**

```bash
curl -X POST http://localhost:8080/profiles \
-H "Content-Type: application/json" \
-d '{
  "id": "user123",
  "age": 28,
  "gender": "F",
  "location": { "lat": 13.7563, "lon": 100.5018 },
  "interests": ["music", "art", "travel"]
}'
```

3. **Seed dummy data:**

```bash
curl http://localhost:8080/seed?count=1000
```

4. **Get all profiles:**

```bash
curl http://localhost:8080/profiles
```

5. **Get matches for user123:**

```bash
curl http://localhost:8080/match/user123
```

---

## 📦 Project Structure

```
/cmd/main.go         → app entrypoint  
/internal/handlers.go → HTTP handlers  
/internal/models.go   → data models  
/internal/store.go    → in-memory storage  
/internal/matcher.go  → matching logic  
/internal/*_test.go   → unit tests  
/Dockerfile           → container setup  
/README.md           → this file
```


