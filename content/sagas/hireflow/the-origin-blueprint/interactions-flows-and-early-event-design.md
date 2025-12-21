---
title: "Interactions, Flows & Early Event Design"
type: "episode"
saga: "HireFlow"
arc: "The Origin Blueprint"
studio: "WastingNoTime Studio"
status: "in progress"
number: 4
summary: "With the service boundaries defined, we explore how they communicate. Direct calls, asynchronous events, and the early hiring flow reveal the first shape of Hireflow’s behavior — resilient, decoupled, and event-driven."
date: "2025-12-28"
tags:
  - distributed-systems
  - event-driven
  - microservices
  - system-design
  - architecture
  - message-bus
---
### Interactions, Flows & Early Event Design

Now that the system’s boundaries exist on the map, another question appears — and it is one of the most important in any distributed architecture:

**How do these services talk to each other without becoming dependent on each other?**

This question is the hinge between _software that scales_ and _software that collapses under its own weight_.

Because services are not just folders separated by network calls.  
They are **small worlds with independent authority**, and the space between them — the communication — is where reliability is either preserved or lost.

So in this episode, we take the first step into that space.

---

#### 1. Two Ways to Interact: Direct Calls vs. Events

In any distributed system, there are only two fundamental coordination patterns:

1. **Direct synchronous calls**

    - predictable

    - transactional feel

    - tightly coupled

    - sensitive to latency and availability

2. **Asynchronous events**

    - decoupled

    - more resilient

    - eventual consistency

    - harder to reason about at first


Neither is “right” or “wrong”.

But each tells a different story about your system.

At Hireflow, events are the main structure.  
Not because it is fashionable, but because the hiring flow is naturally **event-driven**:

- a job is created

- a candidate applies

- an application enters a new stage

- a notification is sent

- search indexes update

- recruiters perform actions

- the system reacts

If we tried to model this with purely synchronous calls, we would recreate a monolith over HTTP — brittle, slow, and over-structured.

Events give the system space to breathe.

---

#### 2. Designing the First Critical Flow: The Happy Path

Milestone 1 is the “Happy Path.”  
The simplest complete scenario a user experiences:

1. A company posts a job.

2. A candidate discovers the job.

3. The candidate applies.

4. Recruiters see the application.

5. Notifications go out.

6. Search updates.

7. The system remains consistent.

Let’s break it down service by service.

---

#### Step 1 — Company Admin Creates a Job

**Actor:** Company Admin  
**Entry:** Gateway → Company-Jobs Service

When a job is created, Company-Jobs becomes the **source of truth**.
But other services also care about this event:

- Search (to index the job)

- Applications (to allow future applications)

- Notifications (in some scenarios)

So Company-Jobs publishes:

```
job.created
```

**Payload (early sketch):**

```json
{
  "jobId": "guid",
  "companyId": "guid",
  "title": "...",
  "createdAt": "timestamp"
}
```

(Details will evolve; for now we sketch only shape, not depth.)

---

#### Step 2 — Candidate Applies

**Actor:** Candidate  
**Entry:** Gateway → Applications Service

Applications is authoritative here.  
But this creates a cascading responsibility:

- Notifications should inform recruiters

- Search must update the “application count” or visibility

- Recruiters rely on fresh application data

So Applications publishes:

```
application.submitted
```

**Payload (early sketch):**

```json
{
  "applicationId": "guid",
  "candidateId": "guid",
  "jobId": "guid",
  "submittedAt": "timestamp"
}
```

This event becomes the heart of the hiring workflow.

---

#### Step 3 — Recruiter Visibility

Recruiters consume the application list from Applications Service via the Gateway.

But they do **not** make direct calls to Candidates Service for candidate profiles.  
That would create dependency loops.

Instead, Applications only stores candidate identifiers.  
If Applications needs candidate display data later, it listens to:

```
candidate.updated
```

or

```
candidate.created
```

and caches what it needs internally.

This ensures that the flow remains stable even if one service is down.

---

#### Step 4 — Notifications Triggered

Notifications Service subscribes to many early events, but one is foundational:

```
application.submitted
```

When it receives this event, it:

- formats notification templates

- routes messages to recruiters

- delivers messages asynchronously

This also allows throttling, retry strategies, and idempotency — all crucial for reliability.

---

#### Step 5 — Search Updates

Search does not call anyone.  
Search listens.

From day one this rule keeps Search simple, reactive, and safe to rebuild.

Events Search listens to:

- `job.created`
- `job.updated`
- `application.submitted`
- (later) `candidate.updated`

Search builds its own denormalized views optimized for querying.

This gives us the first mental diagram:

```
           job.created
       ┌──────────┴───────────┐
Company-Jobs             Search Service
       │
       └──────> application.submitted
                      │
                      ├──> Notifications
                      └──> Search Service (again)
```

This is not the final design, but it captures the spirit:
**events radiating outward, each service reacting according to its responsibility.**

---

#### 3. Event Naming Conventions (Early Version)

Simple. Verb past tense. Domain-prefixed.

Examples:

- `job.created`

- `job.updated`

- `application.submitted`

- `application.status_changed`

- `candidate.created`

- `candidate.updated`

Clear, predictable, semantically stable.

We avoid things like `jobCreateEventV2`  
because events should evolve by payload, not name.

---

#### 4. Event Contracts: The Early Philosophy

At this early stage, our contracts follow three rules:

##### Rule 1 — Minimal but Sufficient

Include only what other services truly need.  
If something can be fetched by ID later, send IDs.

##### Rule 2 — Never Break Consumers

If an event must evolve, extend the payload — do not change meaning.

##### Rule 3 — Events Describe Facts, Not Commands

Events are past tense.  
They announce what happened, not what should happen.

This keeps publishers honest and consumers stable.

---

#### 5. Fault Tolerance at the Edges

One of the early questions we must answer:

**“What happens when a service is down?”**

Event-driven architecture answers this elegantly:

- Services publish events independently.

- Consumers pick them up when they come back online.

- The system self-heals through message persistence.

This is why Hireflow chooses RabbitMQ early — not for scale, but for **predictable behavior during partial failures**.

It’s not perfect, but it’s reliable.

---

#### 6. The First Architecture Flow

Here’s the distilled early picture — the skeleton that Milestone 1 will bring to life:

```
Gateway
   │
   ▼
Identity  ←→  Company-Jobs  → emits job.created
                        │
                        └─────→ Search (index job)

Gateway
   │
   ▼
Applications → emits application.submitted
       │
       ├────→ Notifications
       └────→ Search (update indexes)
```

A small system, already alive.

Small enough to reason about.
Large enough to reveal its nature.

---

#### Closing Reflection

Interaction design is not about drawing arrows.  
It is about understanding _who depends on what_  
and _how to protect each service from accidental entanglement_.

In Hireflow, events become a quiet contract between domains.

They give us:

- resilience

- evolution

- isolation

- clarity

…which are the ingredients that allow a small system to become a reliable one.

Episode 5 will take us deeper into the specifics:

- modeling the first real events

- shaping the topic structure

- choosing idempotency strategies

- designing the “Happy Path” sequence concretely

