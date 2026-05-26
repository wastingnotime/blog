---
name: repository-contract-boundary
description: Define or evolve repository boundary contracts, public surfaces, and visitor guidance. Use when creating a new repository, reworking AGENTS.md, exporting contracts, or reducing accidental coupling between repositories.
---

# Repository Contract Boundary

This repository owns the boundary-contract skill so repository identity stays explicit as the overlay evolves.

## Purpose

Use this skill to define a repository as a bounded semantic actor:

- what it owns
- what it exports
- what remains internal
- how visitors should navigate it safely
- how it participates in campaigns, findings, and integrations

The goal is semantic locality, not repository sprawl.

## When To Use

Use this skill when:

- creating a new repository
- revising repository boundaries
- rewriting or tightening `AGENTS.md`
- defining exported contracts or public folders
- introducing or clarifying `/contracts`
- reducing accidental coupling or semantic leakage
- preparing a repository for repeated AI visitation
- formalizing campaign, findings, or integration expectations

## Boundary Model

Define the repository in this order:

1. Repository role
2. Responsibilities
3. Non-responsibilities
4. Exported semantics
5. Stable vs experimental vs internal surfaces
6. Visitor guidance
7. Participation semantics

## What To Produce

### 1. Repository Identity

State clearly:

- what the repository owns
- what it does not own
- why it exists
- who should visit it
- what other repositories may rely on

### 2. AGENTS.md Rework

Rework `AGENTS.md` into visitor-facing boundary guidance.

It should route readers toward:

- repository identity
- responsibilities
- non-responsibilities
- public contracts
- stable vs experimental areas
- internal-only areas
- campaign / findings interaction
- ownership expectations
- safe navigation paths

Keep it short and explicit. Visitors should not have to infer the boundary.

### 3. Public Contract Map

Propose an appropriate `/contracts` layout for the repository type, such as:

```text
/contracts/api
/contracts/runtime
/contracts/auth
/contracts/events
/contracts/telemetry
/contracts/integration
/contracts/versioning
```

Use only the contract folders that matter for the repository. The point is to expose stable semantics, not mirror implementation internals.

### 4. Export Map

Classify exposed surfaces into:

- `Stable`: safe long-term exported semantics
- `Experimental`: evolving or unstable areas
- `Internal Only`: implementation details that should not be depended upon

### 5. Visitor Guidance

Explain how a visiting agent should approach the repository:

- read the boundary doc first
- prefer public surfaces over internals
- assume exported semantics are the contract
- inspect internals only when the public surface is insufficient
- avoid coupling to implementation details
- use campaigns or findings for cross-repo coordination when needed

## Principles

### Repository Diplomacy

The repository should say how it wants to be approached. Do not leave visitors guessing.

### Semantic Firewall

Visitors should depend on exported semantics, not on private structure.

### AI-Oriented Navigation

Optimize for:

- bounded cognition
- semantic locality
- token efficiency
- reasoning stability
- low-entropy traversal

### Avoid Overreach

Do not turn one repository into the place where all organizational knowledge accumulates.

## Output Shape

When applying this skill, produce concise repository guidance that a transient agent can use immediately.

The expected end state is:

1. a clear repository identity
2. a rewritten `AGENTS.md`
3. a public contract map where relevant
4. stable / experimental / internal separation
5. explicit visitor and participation guidance

## Relationship To Campaigns And Findings

Use boundary contracts to coordinate safely across repositories.

- campaigns coordinate change
- findings record durable observations
- contracts define the safe surface of interaction

Repositories should evolve independently while staying explicit about what they export.

