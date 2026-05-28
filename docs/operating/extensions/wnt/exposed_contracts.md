# WNT Exposed Contracts Extension

## Purpose

This document extends the consuming repository's MRL core `docs/operating/exposed_contracts.md` guidance for WNT-managed repositories.

The core exposed-contracts guidance stays portable. This extension captures the WNT-specific operating rule for repositories whose released artifacts are consumed by other WNT repositories, infrastructure, operators, or external workflows.

## Scope

This document is WNT overlay guidance.
It does not replace MRL core behavior and does not make every repository publish a formal API contract.

Use it when a WNT repository exposes an API, event stream, webhook, CLI surface, operational interface, or deploy parameter contract that another repository or operator is expected to rely on.

---

## Extension Rule

When a WNT repository has an externally relied-on interface, that interface should be treated as an exposed contract.

The repository should:

- keep the contract artifact in a tracked, stable location
- review meaningful contract changes during `release`
- include contract references in release notes when behavior changes
- include contract references in exposure evidence when the released artifact is handed off or published
- classify meaningful changes as breaking, additive, or internal-only

This extends the starter guidance by making the expectation concrete for WNT repositories that participate in cross-repository workflows.

---

## WNT Contract Locations

Prefer these locations when they fit:

```text
docs/contracts/api/
docs/contracts/events/
docs/contracts/ops/
docs/contracts/integration/
```

Use concise Markdown summaries when a formal OpenAPI, AsyncAPI, JSON Schema, Protobuf, or equivalent artifact is not justified yet.

---

## Release Evidence

When a request changes an exposed contract, release evidence should record:

- contract artifact path
- compatibility classification
- affected consumers or downstream repositories when known
- release notes reference
- integration summary reference when ecosystem validation applies

Do not treat an infrastructure PR, issue, or deployment manifest as the contract itself. Those artifacts may point to the contract, but they do not replace it.

---

## Exposure Evidence

When the release is exposed through a WNT lifecycle path, `work/changes/<id>/exposure.md` should record:

- released contract artifact paths or URLs
- whether the exposed contract changed
- compatibility classification
- publication, handoff, or deployment evidence that points to the contract
- any consumer, operator, or downstream feedback that should become evidence for a later `extract` pass

For the WNT AWS ECR, integration validation, and infra promotion path, use this extension alongside `docs/operating/extensions/wnt/expose_aws_ecr_infra_pr.md`.

---

## Anti-Patterns

- claiming contract compatibility from code shape alone
- treating a deployment manifest as a consumer contract
- publishing a runtime artifact while leaving consumers to reverse engineer its interface
- changing a relied-on payload or API without release evidence
- routing consumer feedback directly into implementation without first recording evidence for the next loop

---

## Related Files

- core guidance in consuming repositories: `docs/operating/exposed_contracts.md`
- WNT exposure guidance: `docs/operating/expose_extensions.md`
- AWS ECR, integration validation, and infra promotion extension: `docs/operating/extensions/wnt/expose_aws_ecr_infra_pr.md`
- WNT release delivery validation: `docs/operating/extensions/wnt/release_delivery_validation.md`
