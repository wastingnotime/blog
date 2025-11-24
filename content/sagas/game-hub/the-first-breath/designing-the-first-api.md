---
title: "Designing the First API — /duel/start"
type: "episode"
saga: "Game Hub"
arc: "Trivia Duel"
studio: "WastingNoTime Studio"
status: "in progress"
number: 4
summary: "The first endpoint of Trivia Duel — /duel/start — defines how intention becomes structure. A simple request opens the door for matchmaking, concurrency, and the quiet choreography behind a duel."
date: "2025-11-24"
tags:
  # Domain
  - trivia-duel
  - matchmaking
  - duel-flow
  # Technology
  - golang
  - backend
  - distributed-systems
  - api-design
  - system-design
  # Topic
  - architecture
  - game-dev
---
### Designing the First API: `/duel/start`

Every system begins with a single request.  
For our Game Hub, that first one is humble but symbolic:

```text
POST /duel/start
```

It’s the handshake between a player’s curiosity and the system’s logic.  
And the funny thing about beginnings is — they look simple, but they set the tone for everything that follows.

---

#### Defining What “Start” Means

When a player calls `/duel/start`, what are they _really_ asking?

Not to “play right now,” but to _declare intention_:

> “Hey system, I’m ready. Pair me with someone.”

That’s not a duel yet — it’s a **challenge request**.

We could easily jump straight into creating a duel record, but that would break scalability.  
Real systems need breathing room — time to match, to handle concurrency, to fail gracefully.  
So instead of forcing the duel into existence, we’ll queue that intention and let another component (the matchmaker) decide when the duel truly begins.

---

#### Step 1 — The Data Model

Let’s start small — in Go, we define a `Challenge` struct.

```go
package domain

import "time"

type Challenge struct {
    ID        string    `json:"id"`
    PlayerID  string    `json:"player_id"`
    Status    string    `json:"status"`     // waiting, matched, cancelled
    CreatedAt time.Time `json:"created_at"`
}
```

Each challenge is like a small heartbeat — alive for a short time, waiting for its counterpart.

Later, we might add things like `SkillLevel`, `Region`, or `GameMode` to improve matchmaking, but for now, minimalism wins.

---

#### Step 2 — The Handler

In Go, handlers are where ideas meet HTTP reality.  
Our first one looks like this:

```go
package handlers

import (
    "encoding/json"
    "net/http"
    "time"
    "github.com/google/uuid"
    "wastingnotime/gamehub/domain"
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

    // TODO: Save challenge (in-memory, Redis, or DB)
    // For now, just return it
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(challenge)
}
```

No database yet, no queues, no fancy logic — just structure.  
You could call this the “paper prototype” of an API: something that looks real, but is light enough to move around without fear.

---

#### Step 3 — The Philosophy of Simple Starts

Why so basic?  
Because complexity grows by itself — you never need to plant it.

At _WastingNoTime_, the rule is:

> Build something that runs,  
> then build something that matters.

The first build is for _you_ — to visualize the flow.  
The second is for _the system_ — to stand on its own.

So before touching Redis, DynamoDB, or even Docker, we’ll test this endpoint locally with a simple request:

```bash
curl -X POST http://localhost:8080/duel/start \
  -H "Content-Type: application/json" \
  -d '{"player_id": "P001"}'
```

Response:

```json
{
  "id": "df92ac8e-4c92-4bb2-932f-98e71c6e4db5",
  "player_id": "P001",
  "status": "waiting",
  "created_at": "2025-10-31T21:00:00Z"
}
```

That’s it.  
A first breath of life.

---

#### Step 4 — Looking Ahead

Now we have our entry point.  
From here, we can:

- Persist this challenge in memory (for testing).

- Create a simple **matchmaker goroutine** that checks for pairs every few seconds.

- Generate a duel once two players are ready.


Step by step, the _hub_ will start to emerge — not from a blueprint, but from dialogue between code and curiosity.

---

#### Closing Thought

Software isn’t born perfect — it grows like a story.  
And every story starts with a scene that seems small, until you look back and realize it was everything.

That’s our `/duel/start`.  
A humble endpoint, but the root of all motion.

---

**Source for this episode:**  
Tag **v0.1.0-e04-trivia-duel**  
<https://github.com/wastingnotime/game-hub/tree/v0.1.0-e04-trivia-duel>
