# Expose Extensions

## Purpose

This document defines how WNT repositories can add reusable post-release exposure mechanisms without redefining MRL core or turning lifecycle behavior into an implementation pack.

## Scope

This document defines optional lifecycle guidance after MRL release.
It does not define MRL core deployment behavior, a mandatory MRL phase, or organization-specific release coordination.

Use it when a repository repeatedly exposes released artifacts through the same operational path and wants that path to be explicit, reusable, and optional.

---

## Why Extensions

MRL release accepts an implemented request state as the intended internal version.
WNT exposure guidance describes what happens after that accepted state needs contact with a real context.

MRL core does not prescribe an exposure phase implementation, deployment platform, CI system, registry, or infrastructure handoff shape.

Implementation packs are also not the right place for this concern. Packs define implementation defaults such as language, architecture, source layout, and runtime topology. Exposure mechanisms are lifecycle defaults layered on top of the selected pack.

An expose extension is the WNT overlay way to capture that lifecycle shape.

---

## Definition

An expose extension is a reusable WNT lifecycle description of how accepted release artifacts are exposed into a recurring target context.

An expose extension may define:

- exposure targets and handoff boundaries
- packaging and publication steps
- deployment or handoff automation
- environment parameter contracts
- evidence and feedback channels
- failure classes specific to the exposure path
- production traceability expectations

An expose extension must not:

- redefine the MRL loop
- replace `release`
- turn post-release lifecycle evidence into a mandatory core phase
- silently change domain semantics
- act as a substitute for long-term operations ownership
- claim production deployment without infrastructure-owned deployment truth

---

## Recommended Location

Store reusable expose extensions under:

```text
docs/operating/extensions/
```

Suggested naming shape:

```text
docs/operating/extensions/expose_<mechanism>.md
```

Examples:

- `docs/operating/extensions/wnt/expose_aws_ecr_infra_pr.md`
- `docs/operating/extensions/expose_github_actions_ecs.md`

---

## Repository Decision Rule

If a repository adopts one expose extension as its default path, record that in `decisions.md`.

That decision should say:

- which expose extension is the default
- what boundary counts as exposure completion
- what artifact contract must exist before exposure starts
- what downstream system or repository receives the handoff

Suggested wording:

```md
## DEC-XXXX - Default Exposure Extension

- Date: YYYY-MM-DD
- Status: accepted
- Owners: both

### Context
This repository repeatedly exposes accepted artifacts through the same remote infrastructure path. The exposure mechanism should be reusable and explicit without coupling MRL core to one platform.

### Decision
This repository adopts `<extension_name>` as its default expose extension.

The extension defines how released artifacts are:

- built into an exposure-ready runtime artifact
- published to `<registry or artifact target>`
- handed off to `<infra repository, deployment controller, or environment>`
- traced through explicit exposure evidence

Exposure completion for this repository means `<for example: opening the infrastructure PR>` rather than silently assuming that merge or production rollout already occurred.

Production truth remains with the infrastructure deployment manifest and immutable artifact digest. Release notes and exposure evidence can point to that truth, but they do not replace it.

### Consequences
Exposure behavior becomes easier to repeat and inspect across slices. The repository must keep the extension contract and its evidence artifacts current when the delivery path changes.

### Alternatives considered
Treat exposure as ad hoc per change, or encode the mechanism as part of a language/runtime pack. These were rejected because the mechanism is lifecycle-specific rather than implementation-specific.
```

---

## Per-Change Exposure Rule

Each actual exposure event should be recorded per change under:

```text
work/changes/<id>/exposure.md
```

The extension describes the default mechanism.
The per-change artifact records what actually happened for one released version after release acceptance.

---

## Minimal Extension Spec

Each expose extension should define at least the following sections.

### 1. Identity

- extension name
- intended use cases
- situations where the extension is a bad fit

### 2. Exposure boundary

- what counts as the exposure target
- what counts as successful handoff
- what remains outside the repository's responsibility

### 3. Required inputs

- release decision artifact
- released version identifier
- runtime artifact contract
- credentials or identity assumptions
- environment or parameter contract
- target repository or environment references
- release notes when the artifact has user-facing, contract, runtime, or operational impact
- integration summary when ecosystem validation applies

### 4. Outputs

- exposure event or exposure plan
- published artifact identifiers such as image tag and digest
- handoff references such as workflow URL or infrastructure PR URL
- initial feedback channels
- feedback signals that should become source evidence for the next `extract` pass
- production manifest reference when rollout is observed
- release notes reference
- integration summary reference
- explicit blockers or mismatches when exposure is partial

### 5. Invariants

- exposure only starts from an accepted released version
- the exposed artifact is portable and runtime-ready
- handoff metadata is explicit and reproducible
- the extension does not silently mutate business behavior during exposure
- completion status is recorded even when the target deploy does not occur
- production state is inferred from infrastructure manifests and immutable digests, not from Git history, release notes, or issues alone

### 6. Failure classes

- packaging failure
- publication failure
- contract mismatch between application and infrastructure
- credential or permission failure
- handoff rejection
- downstream deploy blocked or deferred

### 7. Evidence sources

- CI or GitHub Actions run
- published artifact digest
- deployment manifest or parameter diff
- infrastructure pull request
- runtime health or smoke signal when available
- release notes
- integration summary

---

## `exposure.md` Schema

Use this as a minimal per-change template when an expose extension is active.

```md
# Exposure

- Change: `<change_id>`
- Date: YYYY-MM-DD
- Released version: `<version>`
- Exposure extension: `<extension_name>`
- Exposure target: `<environment, workflow, or handoff target>`
- Exposure boundary: `<what counts as done for this repository>`
- Status: planned | partial | completed | blocked

## Inputs

- Release decision: `work/changes/<id>/release_decision.md`
- Artifact reference: `<image tag, digest, bundle, or package>`
- Runtime contract: `<env bundle, params file, manifest, or equivalent>`
- Release notes: `work/changes/<id>/release_notes.md`
- Integration summary: `work/changes/<id>/integration_summary.md` or `<external summary reference>`

## Execution

1. `<step actually performed>`
2. `<step actually performed>`
3. `<step actually performed>`

## Evidence

- Workflow run: `<url or id>`
- Published artifact: `<uri>`
- Digest: `<digest>`
- Handoff artifact: `<infra PR url, manifest diff, etc>`
- Deployment manifest: `<production manifest path, params path, or n/a>`
- Production digest observed: `<digest or n/a>`
- Runtime signal: `<health check, smoke output, or n/a>`

## Result

Describe what exposure achieved in concrete terms.

Do not claim production deployment unless the production deployment manifest or environment state confirms the immutable artifact digest.

## Gaps And Mismatches

- `<contract mismatch>`
- `<blocked deployment step>`
- `<missing downstream confirmation>`
- `<missing release notes, integration summary, or production digest trace>`

## Follow-Up Evidence

- record new runtime, stakeholder, operator, or consumer signals as source evidence for the next `extract` pass
- return to `extract` when the exposure reveals new facts, friction, surprises, drift, or feedback
- return to `refine` only after extracted evidence changes the model or request boundary
- return to `build` only after refinement identifies an implementation slice
```

---

## Feedback Routing

Feedback is not a shortcut around the loop.
Signals created by exposure should be captured as source evidence, then routed into the next `extract` pass.

Examples include:

- operator notes from rollout or handoff
- consumer reports after contract exposure
- runtime health, smoke, alert, metric, trace, or log evidence
- downstream repository comments on an infrastructure PR
- stakeholder feedback after a released workflow is used

Do not patch behavior directly from exposure feedback unless the repository has first recorded the evidence and decided the appropriate loop entry point.

---

## Reuse Rule

Promote an expose extension only after it has been used across real slices or repositories and the repeated shape is clear.

When promoting one:

1. keep MRL core unchanged
2. keep implementation pack guidance separate
3. isolate platform-specific exposure assumptions in the extension document
4. keep per-change exposure evidence in `work/changes/<id>/exposure.md`
5. record repository adoption in `decisions.md`
6. keep release notes, integration summaries, deployment manifests, campaigns, blockers, and findings semantically separate

---

## Related Files

- generic guidance: `docs/operating/expose_extensions.md`
- WNT concrete extension: `docs/operating/extensions/wnt/expose_aws_ecr_infra_pr.md`
- WNT release delivery guidance: `docs/operating/extensions/wnt/release_delivery_validation.md`
- reusable per-change template: `work/changes/_template/exposure.md`
