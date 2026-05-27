# WNT Release Delivery And Validation

## Purpose

This guidance defines how WNT repositories should record delivery facts so operators and agents can answer:

- what is deployed in production
- what changed
- what was validated
- why the change exists
- which dependencies or anomalies affected delivery

Each artifact has one responsibility. Do not overload issues, deployment manifests, release notes, or validation reports with meanings owned by another artifact.

This document is WNT extension guidance. It explains how to trace delivery facts across the repositories that own them.

## Authority

Delivery validation authority is split by repository:

- producing repository: release decision, release notes, candidate image metadata, exposed contract references, and exposure evidence
- `integration-sandbox`: candidate validation, validation run evidence, validation result, and promotion handoff evidence
- `infra-platform`: ECR registry authority, Swarm runtime, deployment manifests, production desired state, and production promotion decision

If validation guidance here conflicts with `integration-sandbox` contracts, `integration-sandbox` wins for validation.
If production guidance here conflicts with `infra-platform` contracts, `infra-platform` wins for registry, deployment, Swarm runtime, and production status.

Read:

- `integration-sandbox/contracts/deployment/candidate-validation-handoff.md`
- `integration-sandbox/contracts/test-scenarios/README.md`
- `integration-sandbox/contracts/runtime/README.md`
- `infra-platform/docs/contracts/deployment/production-promotion-authority.md`
- `infra-platform/docs/contracts/integration-guidance.md`
- `infra-platform/docs/contracts/swarm/swarm-runtime-contract.md`

---

## Artifact Responsibilities

| Question | Source of truth |
|---|---|
| What is deployed in production? | infrastructure deployment manifests plus immutable image digest |
| What changed? | release notes |
| Was it validated in the ecosystem? | integration-sandbox summary |
| Why was the ecosystem changed? | campaign |
| What urgent dependency blocked progress? | blocker |
| What anomaly or mismatch was observed? | finding |

---

## Production Truth

The production source of truth is the infrastructure layer that controls rollout.

For containerized services this means:

- infrastructure repository pull requests
- deployment manifests, parameters, or environment configuration
- immutable image digest or equivalent artifact identifier

Application repositories should not claim that a change is in production only because code was merged, a release decision was accepted, an image was published, or an infra PR was opened.

To know what is in production:

1. inspect the production deployment manifest or deploy parameter state
2. read the image URI and digest, not only the mutable tag
3. map the image digest back to the source revision and release notes
4. confirm whether the target environment rolled out that manifest

---

## Image Metadata

Container images should carry OCI-compatible provenance labels when the build pipeline supports them.

Required baseline labels:

```dockerfile
LABEL org.opencontainers.image.source="<github repository url>"
LABEL org.opencontainers.image.revision="$GIT_SHA"
LABEL org.opencontainers.image.created="$BUILD_DATE"
```

Recommended WNT labels:

```dockerfile
LABEL wnt.integration-summary="<integration summary id or url>"
LABEL wnt.release-notes="<release notes path or url>"
LABEL wnt.campaign="<campaign id or url>"
LABEL wnt.build-url="<ci workflow url>"
```

If a repository cannot add one of these labels yet, record the gap in release or exposure evidence instead of inventing a placeholder.

---

## Release Notes

Release notes explain what changed in the delivered artifact.

They should include:

- delivered features
- fixes
- infrastructure changes
- breaking changes
- migrations
- observability changes
- related campaigns
- resolved blockers
- validation summary references
- image digest or package identifier when available

GitHub Issues are not release notes.

Issues describe work, discussion, and task state.
Release notes describe delivered behavior, delivered capability, and operational impact.

Use `work/changes/_template/release_notes.md` for per-change release notes when the repository does not already have a stronger local release-note convention.

---

## Integration Sandbox Summary

Integration summaries describe ecosystem validation.

They should include:

- validated repositories
- validated versions or image digests
- executed scenarios
- result status
- logs, traces, screenshots, or other evidence

The integration summary answers whether a version was validated. It does not decide deployment, replace release notes, or become campaign memory.

Use `work/changes/_template/integration_summary.md` as the baseline shape when recording local validation summaries or when preparing a request for `integration-sandbox`.

For production-bound container images, the authoritative validation output is the validation run and promotion handoff evidence emitted by `integration-sandbox`.

---

## Campaigns, Blockers, And Findings

Campaigns explain why coordinated ecosystem change exists.

Blockers identify urgent external dependencies that prevent safe progress.

Findings document anomalies, mismatches, inconsistencies, or unexpected behavior.

Keep these artifacts separate:

- do not use a campaign as a release note
- do not use a blocker for non-blocking observations
- do not use a finding to imply ownership transfer or urgency
- do not use integration evidence as long-term coordination memory

---

## Delivery Flow

Use this default flow when a WNT repository delivers a runtime artifact:

```text
Repository Change
    -> CI Build
    -> Integration Sandbox Validation
    -> Validation Result And Promotion Handoff
    -> Release Notes
    -> Infra-Platform Promotion Decision
    -> Deployment Manifest Update When Accepted
    -> Rollout Observation
```

The exact automation may differ by repository, but the artifact responsibilities should stay stable.

---

## Repository Responsibilities

| Repository type | Responsibility |
|---|---|
| application repository | build artifacts, release notes, image metadata, exposure evidence |
| integration-sandbox | ecosystem validation runs, validation results, and promotion handoff evidence |
| infra-platform | registry authority, Swarm runtime, deployment truth, and production manifest state |
| management | campaign coordination and durable cross-repository memory |
| observing repository | findings |
| blocked repository | blockers |

---

## Agent Rule

When asked what changes are in production, do not answer from Git history or release notes alone.

Start from the infra-platform production manifest state and immutable image digest, then trace backward to:

- source revision
- release notes
- integration summary
- campaign, blocker, or finding references when present

If any link is missing, say which artifact is missing and treat the production answer as partial.
