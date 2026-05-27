# WNT Expose Extension: AWS ECR, Integration Validation, Infra Promotion

## Purpose

Use this extension when a WNT repository exposes an accepted release as a container image that must be validated by `integration-sandbox` before `infra-platform` promotes it to production infrastructure.

This is a coordination path, not a repository-local deployment path.

## Scope

This document is WNT overlay guidance.
It does not define MRL core behavior, grant production authority to the producing repository, or replace the contracts owned by `integration-sandbox` and `infra-platform`.

The authoritative repositories are:

- producing repository: owns source behavior, release decision, image build, release notes, and candidate evidence
- `integration-sandbox`: owns candidate validation and promotion handoff evidence
- `infra-platform`: owns the container registry authority, Swarm production runtime, deployment manifests, and production promotion decision

---

## Good Fit

Use this extension when:

- the released artifact is a container image
- the image is intended for WNT production infrastructure
- `integration-sandbox` has or should have a validation contract for the candidate
- `infra-platform` owns the ECR repository, Swarm cluster, deployment manifests, and promotion authority
- the producing repository needs to request guidance or handoff rather than directly deploy

Do not use this extension when:

- the repository deploys directly under its own production authority
- the artifact is not a container image
- validation does not require `integration-sandbox`
- production does not run through `infra-platform`

---

## Authority Boundary

The producing repository may:

- accept the release internally
- build or request publication of a candidate image according to the infra-owned registry contract
- provide release notes, image metadata, contract references, and validation inputs
- request validation from `integration-sandbox`
- request promotion guidance from `infra-platform`
- record exposure evidence and gaps

The producing repository must not:

- claim that an image is production-ready without `integration-sandbox` validation evidence
- claim production deployment without `infra-platform` deployment truth
- mutate infra deployment manifests directly unless `infra-platform` guidance explicitly delegates that action
- treat image publication or validation as production promotion

`integration-sandbox` owns whether the candidate image is validated for promotion.

`infra-platform` owns:

- ECR repository and registry authority
- Swarm cluster and runtime placement
- deployment manifests and production desired state
- final production promotion or rejection

---

## Read First

Before recording exposure, read the owning repository guidance:

- `integration-sandbox/contracts/deployment/README.md`
- `integration-sandbox/contracts/test-scenarios/README.md`
- `integration-sandbox/contracts/runtime/README.md`
- `infra-platform/docs/contracts/integration-guidance.md`
- `infra-platform/docs/contracts/deployment/dispatch-contract.md`
- `infra-platform/docs/contracts/swarm/swarm-runtime-contract.md`

If either owner repository lacks the needed contract for the candidate, record the gap and open a finding or campaign instead of inventing local deployment procedure.

---

## Required Inputs

- accepted `work/changes/<id>/release_decision.md`
- release notes for behavior, contracts, runtime impact, and operator-visible changes
- integration summary or validation request context
- image name, tag, digest, repository, and build provenance when available
- contract references for APIs, events, runtime expectations, health endpoints, telemetry, and deploy parameters
- candidate manifest or equivalent validation input expected by `integration-sandbox`
- infra guidance or issue/PR from `infra-platform` when promotion is requested

---

## Execution Shape

1. Confirm that the request state is accepted by `release`.
2. Prepare release notes and any exposed contract references.
3. Produce or request the candidate image according to the infra-owned registry guidance.
4. Ask `integration-sandbox` for the current candidate validation path.
5. Provide the candidate manifest, image digest, release notes, contract references, and seed/runtime assumptions requested by `integration-sandbox`.
6. Wait for or link validation evidence from `integration-sandbox`.
7. If validation passes, ask `infra-platform` for promotion guidance or handoff handling.
8. Record the infra-owned deployment manifest, PR, issue, or rejection evidence.
9. Record whether the exposure boundary reached candidate publication, validation handoff, validation pass, infra handoff, manifest promotion, or observed rollout.
10. Route feedback, runtime surprises, validation failures, or infra blockers into evidence for a later `extract` pass.

---

## Evidence To Record

In `work/changes/<id>/exposure.md`, record:

- exposure extension: `wnt/aws_ecr_infra_validation`
- producing repository and released version
- image URI, tag, and immutable digest when available
- candidate manifest or validation request reference
- `integration-sandbox` validation run, result, and handoff artifact
- `infra-platform` issue, PR, deployment manifest, or promotion decision
- production status:
  - not requested
  - validation requested
  - validation failed
  - validation passed
  - infra handoff requested
  - infra accepted
  - production rollout observed
  - infra rejected or blocked
- release notes and contract references
- gaps, blockers, and ownership questions

---

## Production Truth

Production status is not inferred from:

- release acceptance
- image build success
- image publication
- integration validation alone
- a producing-repository issue or exposure note

Production status is confirmed only by `infra-platform` deployment manifest or environment evidence that references the immutable image digest.

---

## Failure Classes

- candidate image missing or digest not immutable
- candidate manifest incomplete
- validation contract missing in `integration-sandbox`
- validation failed or inconclusive
- infra registry or Swarm contract missing in `infra-platform`
- infra promotion rejected, blocked, or deferred
- deployment manifest does not reference the validated digest
- rollout cannot be observed

---

## Repository Decision Template

```md
## DEC-XXXX - Default WNT Exposure Extension

- Date: YYYY-MM-DD
- Status: accepted
- Owners: both

### Context
This repository exposes accepted releases as container images that require WNT validation and infra promotion before production.

### Decision
This repository adopts `wnt/aws_ecr_infra_validation` as its default exposure extension.

The producing repository owns release evidence and candidate image metadata.
`integration-sandbox` owns validation evidence.
`infra-platform` owns registry, Swarm runtime, deployment manifests, and production promotion authority.

Exposure completion for this repository means `<candidate publication | validation request | validation pass | infra handoff | production rollout observed>`.

### Consequences
The repository records exposure evidence without claiming authority owned by integration or infrastructure repositories.
```

---

## Related Files

- WNT exposed contracts: `docs/operating/extensions/wnt/exposed_contracts.md`
- WNT release delivery validation: `docs/operating/extensions/wnt/release_delivery_validation.md`
- per-change template: `work/changes/_template/exposure.md`
