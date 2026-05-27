# WNT Skills Workflow Extension

## Purpose

This document extends `docs/operating/skills_workflow.md` for WNT-managed repositories.

The core skills workflow stays portable. This extension captures WNT-specific lifecycle skill guidance after release acceptance.

## Scope

This document is WNT overlay guidance.
It does not redefine the MRL core loop.

Use it with:

- `docs/operating/skills_workflow.md`
- `docs/operating/expose_extensions.md`
- `docs/operating/extensions/wnt/expose_aws_ecr_infra_pr.md`
- `docs/operating/extensions/wnt/release_delivery_validation.md`

---

## Repository Skills

WNT repositories may receive lifecycle and support skills from the overlay under `.agents/skills/`.

WNT-owned lifecycle skills should stay overlay-owned when they contain concrete delivery mechanics.
They should not be added to the starter as canonical MRL phases.

```text
.agents/skills/
  expose/
    SKILL.md
  campaign-management-protocol/
    SKILL.md
  cross-repo-findings/
    SKILL.md
  decommission-repository/
    SKILL.md
  local-vault-access/
    SKILL.md
  repository-contract-boundary/
    SKILL.md
```

The overlay removes the old standalone `feedback` skill from targets.
Runtime feedback should be preserved as source evidence and brought back through `extract`.

---

## WNT Lifecycle Guidance After Release

Exposure and feedback are lifecycle concerns that happen after release acceptance.
They may be supported by repository skills and overlay documents, but concrete publication, deployment, handoff, and feedback-routing behavior is not MRL core.

### Expose

**Goal**

- put the released state into contact with a real context
- expose it to users, stakeholders, workflows, or environments

**Input**

- released version
- exposure plan or target context
- exposure-ready artifact, normally packaged as a container image unless another portable runtime artifact is explicitly justified
- release notes for the delivered artifact when applicable
- integration summary when ecosystem validation applies

**Output**

- exposure event
- initial real-world feedback signals
- handoff evidence such as artifact digest, infrastructure PR, deployment manifest reference, and rollout observation when available

**Must not**

- be treated as long-term operation
- be confused with release itself
- assume source code alone is the exposure artifact when a portable runtime artifact should exist
- claim production deployment without infrastructure-owned manifest or environment evidence

### Feedback Evidence

**Goal**

- treat the exposed state as something that now exists in reality
- collect source evidence from friction, surprises, drift, and feedback
- feed those signals into the next `extract` pass

**Input**

- exposed version
- real-world signals

**Output**

- new source evidence for the next `extract` pass

**Must not**

- be confused with general operations ownership
- bypass the loop by making undocumented corrections

---

## Mental Model

```text
extract
  -> refine/build loop
  -> egd
  -> release

WNT lifecycle after release:

release
  -> exposure guidance
  -> source evidence from feedback
  -> extract
```
