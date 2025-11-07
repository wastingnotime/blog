---
title: "Epic Continues — The Duel Begins"
type: "episode"
saga: "Game Hub"
arc: "The First Breath"
studio: "WastingNoTime Studio"
status: "in progress"
number: 2
summary: "The spark of Trivia Duel marks the moment when play meets architecture. Beneath the simplicity of two minds facing a question, patterns of state, timing, and fairness begin to unfold — the quiet architecture behind every duel."
date: "2025-11-07"
tags:
- game-dev
- distributed-systems
- architecture
- reflection
---

### Epic continues — the duel begins

Alright, so we have our north star: **Trivia Duel** — a simple, fun, and competitive game.

The name already brings movement: two minds, one question, and a bit of pride on the line. It’s light enough to test fast, but structured enough to teach us about user sessions, score tracking, concurrency, and fairness — the kind of things that appear in real systems all the time.

But before we start coding anything, let’s slow down.  
We need to understand what _“a duel”_ really means in system terms.

Two players → a match.  
Each match → a sequence of questions.  
Each question → a small interaction cycle.  
Each cycle → requests, responses, and timeouts.

It sounds simple, but inside those lines hide every challenge we love to face: state management, latency, synchronization, and the eternal dance between **frontend experience** and **backend reliability**.

So yes, it’s a “simple game,” but in disguise, it’s a **distributed system** in miniature — perfect terrain for what _wasting no time_ is all about.

---

### Starting by the invisible part

We’ll begin, as usual, by the part no one sees: the backend.  
Not because we don’t care about visuals, but because **structure creates freedom**. Once the rules are clear, everything else can grow on solid ground.

The first steps will be small:

1. Define how players join a duel.

2. Define how questions are delivered and answered.

3. Track progress, score, and results.

4. Keep everything stateless — or as stateless as possible — so it can scale later.


Simple goals, but each one opens doors to concepts like event-driven flows, in-memory caching, and API orchestration.

---

### Designing from curiosity, not from fear

We won’t rush.  
Every design decision will be a conversation: between code, context, and curiosity.

Sometimes the answer will come fast.  
Sometimes, we’ll stop, think, and say, _“This doesn’t feel right yet.”_  
And that’s fine. Because here, there’s no rush to impress — only a steady pace toward understanding.

This is not just a tutorial; it’s a record of thought in motion — one that others can read, fork, and build upon.

So, let’s prepare the ground for our duel.  
Next: we’ll describe **how a duel is born** — from a simple API call to the first spark of competition.

