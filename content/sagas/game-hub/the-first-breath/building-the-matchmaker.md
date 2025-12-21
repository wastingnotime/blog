---
title: "Building the Matchmaker: Where Waiting Meets Destiny"
type: "episode"
saga: "Game Hub"
arc: "the-first-breath"
studio: "WastingNoTime Studio"
status: "in progress"
number: 5
date: "2025-12-03"
summary: "A quiet background loop becomes the first true heartbeat of the Game Hub, pairing waiting challengers into living duels and revealing the beauty of coordination, state, and synchronization."
tags:
  # Domain
  - game-hub
  - trivia-duel
  - matchmaking
  # Technology
  - golang
  - backend
  - distributed-systems
  - system-design
  # Topic
  - architecture
  - session-flow
  - creation
  # Process / Meta
  - storytelling
  - reflection
---
### Building the Matchmaker: Where Waiting Meets Destiny

Every world needs a force that connects its wandering souls.  
In our Game Hub, that force is the **matchmaker** — the quiet engine that watches, waits, and says,

> “You two. You’ll duel.”

It’s not glamorous code. It doesn’t print shiny UI or speak to the player.  
But it gives shape to chaos — and that’s what makes it beautiful.

---

#### The Role of the Matchmaker

From the previous episode, we know that every time a player hits:

```text
POST /duel/start
```

we create a `Challenge` and store it with a status of `"waiting"`.  
Now we need a background process that checks for pairs of waiting players and creates a **duel** when two are found.

Sounds simple, right?  
But even a tiny piece of coordination teaches us a lot about concurrency, state, and synchronization — three words that every backend developer eventually learns to respect.

---

#### Step 1 — A Minimal In-Memory World

Let’s begin with something almost poetic: a small in-memory room where challenges wait.

```go
package memory

import "github.com/wastingnotime/game-hub/domain"

var WaitingChallenges = make([]domain.Challenge, 0)
```

It’s fragile, temporary, and perfect for the moment — like a prototype carved in sand before we cast it in stone.

Later, we’ll move this to Redis or DynamoDB, but for now, this is where stories start.

---

#### Step 2 — The Matchmaker Loop

Now, let’s write the heartbeat.  
A loop that wakes up every few seconds, looks at the waiting list, and pairs the first two lonely challengers it finds.

```go
package services

import (
    "fmt"
    "sync"
    "time"
	"github.com/wastingnotime/game-hub/domain"
	"github.com/wastingnotime/game-hub/memory"
)

var mu sync.Mutex

func StartMatchmaker() {
    go func() {
        for {
            time.Sleep(2 * time.Second) // small heartbeat
            matchPlayers()
        }
    }()
}

func matchPlayers() {
    mu.Lock()
    defer mu.Unlock()

    for len(memory.WaitingChallenges) >= 2 {
        p1 := memory.WaitingChallenges[0]
        p2 := memory.WaitingChallenges[1]
        memory.WaitingChallenges = memory.WaitingChallenges[2:]

        duel := domain.NewDuel(p1.PlayerID, p2.PlayerID)
        fmt.Printf("✨ New duel started! %s vs %s → ID: %s\n", p1.PlayerID, p2.PlayerID, duel.ID)
    }
}
```

Nothing fancy — just Go doing what Go does best: lightweight concurrency with almost no ceremony.  
Every few seconds, it looks around, finds a pair, and says, “Let’s go.”

You can feel the system breathing.

---

#### Step 3 — The Duel Entity

Now we need to define what a _duel_ actually is — the thing born from that matchmaker spark.

```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type Duel struct {
    ID        string    `json:"id"`
    Players   []string  `json:"players"`
    State     string    `json:"state"` // starting, active, finished
    CreatedAt time.Time `json:"created_at"`
}

func NewDuel(p1, p2 string) Duel {
    return Duel{
        ID:        uuid.New().String(),
        Players:   []string{p1, p2},
        State:     "starting",
        CreatedAt: time.Now().UTC(),
    }
}
```

A duel is just a record of a connection — ephemeral yet real.  
Each one marks the moment two intentions found each other.

---

#### Step 4 — Plugging It Together

In your `handlers.go`:
```go
package handlers

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/wastingnotime/game-hub/domain"
	"github.com/wastingnotime/game-hub/memory"

	"net/http"
	"time"
)

func StartDuelHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		PlayerID string `json:"player_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	challenge := domain.Challenge{
		ID:        uuid.New().String(),
		PlayerID:  req.PlayerID,
		Status:    "waiting",
		CreatedAt: time.Now().UTC(),
	}

	memory.WaitingChallenges = append(memory.WaitingChallenges, challenge)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}
```


In your `main.go`:
```go
func main() {
    fmt.Println("🕹️ Game Hub starting...")
    services.StartMatchmaker()

    //POST /duel/start
    http.HandleFunc("/duel/start", handlers.StartDuelHandler)

    s := &http.Server{
        Addr: ":8080",
    }

    log.Fatal(s.ListenAndServe())
}
```


Now run it.  
Then, from another terminal, call:

```bash
curl -X POST http://localhost:8080/duel/start -d '{"player_id":"P001"}' -H "Content-Type: application/json"
curl -X POST http://localhost:8080/duel/start -d '{"player_id":"P002"}' -H "Content-Type: application/json"
```

And a few seconds later, in your logs:

```
✨ New duel started! P001 vs P002 → ID: a82b0f32-7f83-4f8b-8f9c-33db76a08a0b
```

It’s alive.

---

#### Step 5 — The First Breath of the Hub

This moment always feels special.  
There’s no frontend, no graphics, no noise — but under the surface, something _connected_.

It’s the invisible part of creation that makes everything else possible.  
From here, we can grow in many directions:

- Persist duels and challenges

- Add timeouts and retries

- Introduce AI opponents

- Expand to WebSockets or event-driven communication


But the essence remains: **two requests met each other and created meaning**.

---

#### Closing Thought

Software is never born whole — it’s assembled from conversations between ideas.  
Here, our matchmaker is the quiet listener in that dialogue.  
It doesn’t decide who wins or loses — it just ensures the right people meet at the right time.

And that, in its own way, is a kind of art.

---

**Source for this episode:**  
Tag **v0.2.0-e05-trivia-duel**  
<https://github.com/wastingnotime/game-hub/tree/v0.2.0-e05-trivia-duel>
