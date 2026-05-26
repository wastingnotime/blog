---
name: decommission-repository
description: Decommission a repository by retiring it in the management ledgers, preserving or archiving its context, and closing out any related campaign or contract records. Use when a repository is being replaced, archived, wound down, or removed from active focus.
---

# Decommission Repository

This repository owns the decommissioning workflow so repository retirement stays explicit and auditable.

Use this skill when a repository should stop being treated as active work and needs a clean retirement path in the management ledger.

## When to use

Use this skill for requests like:

- retire a repository
- archive a legacy repository
- mark a repository as decommissioning or decommissioned
- remove a repository from the main pipeline
- preserve a replaced repository as reference

If the change spans more than one repository, use the campaign-management protocol alongside this skill.

## Workflow

1. Confirm the end state.
   - `decommissioning` means retirement is in progress.
   - `decommissioned` means retirement is complete.
   - `archive` means the repo is old, replaced, or non-meaningful in the GitHub architecture.
2. Load the management sources of truth.
   - `docs/pipeline.md`
   - `docs/sleepy-projects.md`
   - `references/github-architecture.md`
   - `docs/event-log.md`
   - `docs/mrl-managed-repository-contracts.md` when the repo signed a managed contract
   - `references/decommission-checklist.md` for the exact execution checklist
3. Update the ledgers first.
   - Move the repo into `decommissioning` or `decommissioned` in `docs/pipeline.md` if it belongs there.
   - Add or update the repo in `docs/sleepy-projects.md` when it is intentionally hidden from the main view.
   - Update `references/github-architecture.md` when visibility, naming, or archive status changes.
   - Record the retirement in `docs/event-log.md` if it is a durable management event.
4. Preserve the right amount of context.
   - Keep a replaced repository as a reference if it still teaches something.
   - Remove it from active views when it no longer belongs in the main pipeline.
   - Delete a repository only when the architecture or explicit instruction calls for deletion.
5. Close related coordination artifacts.
   - Update any campaign intent, status, decisions, and final summary.
   - Archive completed campaign bundles under `campaigns/archive/` when that structure exists.
   - Close or update repository issues tied to the retirement.
6. Verify the final state.
   - The active pipeline no longer treats the repo as live work.
   - The sleep/archive ledgers reflect the new status.
   - Any managed-contract or campaign records are consistent.

## Reference

Read [`references/decommission-checklist.md`](../../references/decommission-checklist.md) when you need the exact execution order or a concise retirement checklist.

