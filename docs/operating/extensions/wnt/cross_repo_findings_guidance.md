# WNT Cross-Repository Findings Handling

## Purpose

Define how a repository should handle **findings discovered locally** that impact other repositories.

Findings and campaigns belong in `#coordination`. Blockers belong in `#blockers`. Runtime incidents belong in `#ops-alerts`.
See `docs/operating/extensions/wnt/inter_repository_coordination_protocol.md` for the full channel model.
The `#ops-alerts` contract is infra-platform owned and is not opened directly by this repository.

This avoids:
- intrusive changes across repositories
- ownership confusion
- hidden coupling
- incorrect “implementation from the outside”

---

## Core Principle

> A repository can observe and report a problem, but only the owning repository should define and implement the fix.

## Authority

`docs/operating/extensions/wnt/inter_repository_coordination_protocol.md` owns the channel model and cross-repository coordination authority.
This document owns the findings and blockers workflow.

If this document and the coordination protocol disagree about channels, urgency, or ownership routing, the coordination protocol wins.

---

## What is a Finding?

A finding is a non-blocking semantic issue owned by another repository that should be solved eventually.

Common examples:

- contract mismatch
- unexpected behavior from another repository
- integration failure
- runtime inconsistency
- deploy/runtime side-effect

A blocker is a blocking semantic issue likely owned by another repository.

Use `blocker` when:

- the current work cannot proceed safely until another repository changes
- the target repository is likely, but not certain
- the opener should create a target issue and let the owning repository confirm or reject responsibility

Use `finding` when:

- the current work can continue
- the issue is outside the current repository's ownership
- the issue should still be tracked until the owning repository resolves it

---

## Decision: Skill or Guidance?

This is the durable repository guidance for the repo-owned `cross-repo-findings` skill.

The skill handles the repeatable workflow. This document keeps the ownership
and reporting rules stable.

Finding issues are Codex-authored. Humans do not file finding issues manually.

---

## Standard Flow

### Step 1 — Create a Finding (Local Repository)

Create:

```
work/findings/<date>-<short-name>/finding.md
```

Example:

```
work/findings/2026-04-29-candidate-image-dispatch/finding.md
```

For a blocker, the local record should clearly identify the blocked activity and the likely target repository.

---

## Finding Document Template

```md
# Finding: <short name>

## Issue Type

- finding | blocker

## Context

Where this was discovered:
- feature
- deploy
- runtime
- test

## Observed Behavior

What actually happened.

Include:
- payloads
- logs
- traces
- screenshots (if relevant)

## Expected Behavior

What should have happened.

Reference:
- contracts
- docs
- assumptions

## Impact

What breaks because of this:
- deploy
- runtime
- user flow
- data integrity

## Suspected Source

Which repository likely owns the issue.

For blockers, include the likely target repository even when ownership is not fully confirmed.

## Evidence

Concrete data:
- logs
- request/response samples
- event payloads

## Suggested Direction (Optional)

High-level direction only.

Do NOT:
- define internal implementation of another repository
- prescribe exact code changes

## Owning Repository

The repository responsible for fixing the issue.

For blockers, if ownership is uncertain but a target repository is likely, open the target issue anyway and let that repository confirm or reject ownership.

## Local Impact

What this repository must do (if anything):
- adapt temporarily
- block deploy or block the current activity
- add validation
```

---

## Step 2 — Escalate to Management (If Needed)

If the finding affects more than one repository:

- create or update a campaign in `management`
- attach this finding as evidence
- create the target repository issue(s) even if the likely owner is not yet confirmed
- link the local issue to each target issue and link each target issue back to the local issue

Do NOT coordinate multi-repo changes locally.

Before escalation, create a GitHub issue in the discovering repository and post
a short advisory note to Discord `#coordination` (`1502714561923383296`) for findings through
`discordctl.py`. Use the original Discord message as the anchor for follow-up replies or a thread when the tooling supports it. The current reference implementation is `infra-platform/.github/workflows/notify-discord-on-issue-open.yml` and `infra-platform/.github/workflows/notify-discord-on-issue-close.yml`.

The issue-open workflow should create the anchor message with the discovery summary, issue link, and anchor ID or link for reuse. The issue-close workflow should reply to that anchor with the conclusion or resolution status.

---

## Step 3 — Create Repository Handoff

For the owning repository:

- create a Change Packet in `management`
- or communicate via direct handoff
- if the issue is a blocker and ownership was only likely, let the owning repository confirm or reject it in their issue thread

The owning repository must:
- interpret the finding
- decide implementation
- update its own knowledge

---

## Step 4 — Local Repository Behavior

The discovering repository may:

- add validation to prevent silent failures
- log better diagnostics
- document the expected contract
- open a GitHub issue in the discovering repository
- open a target issue in the likely owning repository for findings and blockers when ownership is sufficiently clear
- open likely target issues in other repositories for blockers when ownership is uncertain
- notify Discord `#coordination` (`1502714561923383296`) for findings and `#blockers` (`1502715067186024628`) for active dependency interruptions
- reply to or thread from the original Discord finding message when the tooling supports it

The discovering repository issue is the runtime source of truth for the finding's lifecycle. When it closes, send a completion message to `#coordination` (`1502714561923383296`) for findings or `#blockers` (`1502715067186024628`) for blockers.

Target repository PR handling is defined in `.agents/references/target_repository_pr_handling.md`. Non-code resolution handling is defined in `.agents/references/target_repository_resolution.md`. False-positive handling is defined in `.agents/references/target_repository_false_positive.md`.

The discovering repository must NOT:

- implement fixes for another repository
- assume ownership of foreign behavior
- embed permanent workarounds without explicit decision

For findings, do not leave the issue local-only if the owning repository is known or can be reasonably inferred.
For blockers, do not wait for perfect ownership certainty before opening the target issue when the likely owner is clear enough to identify.

## Source Repository Lifecycle

The discovering repository should close the main finding or blocker issue only when:

- all linked target issues are closed
- any related campaign is closed
- the solution or conclusion has been recorded

If a linked target issue is still open, keep the source issue open.

## Target Repository Lifecycle

The target repository should:

- assign the issue to `hriccio` as soon as possible
- open a PR when code changes are required
- link the issue to the PR
- include the closing keyword in the PR text so the issue closes automatically on merge
- when no code changes are required, write the conclusion and close the issue if solved
- when the issue still needs more work, comment that a revision is needed and keep it open
- when the issue is a false positive or there is nothing this repository can do, comment and close the issue

## Discord Threading

When the Discord tooling supports message lookup and replies:

- post the opening finding or blocker message once as the anchor
- reply to that message for updates, ownership confirmations, revision requests, and completions
- use the same anchor for the target-repository side of the conversation when possible
- record the anchor in the local finding record so it can be reused for every follow-up
- the issue-open notification creates the anchor
- the issue-close notification replies to the anchor with the conclusion or resolution status

If replies or threads are not available, include the original message link in each follow-up so the conversation remains connected.

---

## Anti-Patterns

### 1. Writing implementation.md for another repository

Bad:

```
infra-platform/work/changes/.../implementation.md
```

Why:

- breaks ownership boundaries
- creates coupling
- makes future maintenance unclear

---

### 2. Silent workaround

Bad:

```
if payload.invalid:
    fix_it_here()
```

Why:

- hides real issue
- spreads incorrect logic

---

### 3. Expanding scope locally

Bad:

```
this-repo-fixes-everything.md
```

Why:

- consumes unrelated work
- hides true ownership
- slows down the real fix

---

## Outcome

This pattern keeps local repositories honest:

- they can report
- they can validate
- they can document
- but they do not own foreign fixes
