---
title: "Identifying Our First Services"
type: "episode"
saga: "HireFlow"
arc: "The Origin Blueprint"
studio: "WastingNoTime Studio"
status: "in progress"
number: 3
summary: "From actors to boundaries — we translate the briefing into concrete service domains. Each service emerges from responsibility, authority, and clarity, forming the first stable shape of Hireflow’s architecture."
date: "2025-12-05"
tags:
  - distributed-systems
  - microservices
  - domain-boundaries
  - backend-architecture
  - system-design
  - event-driven
---

#### Identifying Our First Services

Now that the actors are mapped, the natural next step is to ask a deceptively simple question:

**“What exists in this world?”**

Actors show _who_ interacts with the system.  
Services define _where_ responsibility lives.

A service is not “a controller” or “a folder” — it is a **boundary of authority**.  
It owns a part of the domain, the data, and the decisions within it.

So we return to the briefing and extract nouns and verbs that hint at areas of responsibility.  
Not implementation details — _domains of meaning_.

From the original briefing, the ecosystem seems to revolve around:

- companies

- job postings

- recruiters

- candidates

- applications

- evaluations

- notifications

- search

- automated workflows


If we squint, these concepts start clustering naturally — as if the system is revealing its own topology.

Let’s process them one by one with WNT’s “calm, deliberate” approach.

---

##### 1. Identity & Roles

We identified earlier that we have four roles:

- Company Admin

- Recruiter

- Candidate

- Timer (system-triggered behavior)


This already defines our first service:

###### Identity Service

**Responsibility:**

- authentication

- authorization

- role management

- issuing tokens

- basic user lifecycle


**Authority:**  
Owns everything related to identity.

**Reasoning:**  
Identity is foundational.  
Every other service depends on knowing _who_ is acting and _what they can do_.

So the Identity Service is not optional — it is the gatekeeper.

---

##### 2. Company & Job Management

Companies create job postings.  
Recruiters manage them.  
Candidates apply to them.

This gives us a second clear boundary:

###### Company-Jobs Service

**Responsibility:**

- company profiles

- job postings

- recruiter associations

- job visibility rules


**Authority:**  
Owns jobs and their lifecycle.

**Why separate it?**  
Because job management is a domain in itself.  
It evolves independently from candidates or evaluations.

---

##### 3. Candidates Service

Candidates are not just entities — they have profiles, histories, resumes, preferences.

###### Candidates Service

**Responsibility:**

- candidate profiles

- resume metadata

- activity history

- candidate-side settings


**Authority:**  
Owns candidate data and its evolution.

This avoids mixing candidate information with job or application information, which is a trap monoliths often fall into.

---

##### 4. Applications Service

The heart of the hiring flow:  
A candidate applies, a recruiter evaluates, a decision is made.

###### Applications Service

**Responsibility:**

- applications

- status transitions

- evaluation records

- scoring or screening logic (initial, not ML)


**Authority:**  
Owns the lifecycle of an application.

This is where business rules grow dense.  
Separating applications keeps the rest of the system clean.

---

##### 5. Notifications Service

Whenever certain events happen — application submitted, recruiter invited — someone needs to be notified.

###### Notifications Service

**Responsibility:**

- sending emails

- preparing message templates

- acting on events


**Authority:**  
Owns outbound communications.

This service listens to events (RabbitMQ in our architecture) and emits notifications.

---

##### 6. Search Service

Across companies, jobs, and candidates, search becomes an important operational capability.

###### Search Service

**Responsibility:**

- indexing data from other services

- exposing search endpoints

- providing “quick access” queries


**Authority:**  
Owns search indexes and denormalized views.

This service _does not store source of truth data_ — it mirrors other domains through events.

---

##### 7. Gateway

Even though it is not a domain service, the system needs a stable entry point.

###### Gateway

**Responsibility:**

- API routing

- request composition

- rate limiting (future)

- authentication pass-through


**Authority:**  
Owns how the outside world enters the platform.

It is the front door.

---

##### 8. Timer / Automated Workflows

Timers trigger cleanup routines, periodic evaluations, or automated transitions.  
These are orchestrated not inside the services but through scheduled messages.

###### Workflow Scheduler

Not exactly a “service”, but a mechanism:

**Responsibility:**

- producing time-triggered messages

- orchestrating routine tasks


**Authority:**  
Owns temporal events.

This could be implemented with Kubernetes CronJobs, or a lightweight internal scheduler.

---

#### Current Map (Early Draft)

By now, the system begins to reveal its natural structure:

- **identity**

- **company-jobs**

- **candidates**

- **applications**

- **notifications**

- **search**

- **gateway**

- **scheduler (cron/timer)**


This list will evolve.  
It always does.  
But it already gives us enough clarity to design paths, boundaries, and interactions.

And more importantly:  
it gives us confidence that each service has a purpose and a clean responsibility.

