---
name: expose
description: Record WNT post-release exposure evidence after a release decision is accepted. Use when an accepted release needs publication, infrastructure handoff, rollout evidence, or explicit exposure-boundary documentation without treating exposure as an MRL core phase or long-term operations ownership.
---

# Mission

Record how an accepted release is put into contact with a real WNT context.

# Read First

- `docs/operating/expose_extensions.md`
- `docs/operating/extensions/wnt/expose_aws_ecr_infra_pr.md` when the repository uses the WNT AWS ECR, integration validation, and infra promotion path
- `docs/operating/extensions/wnt/release_delivery_validation.md`
- `work/changes/<id>/release_decision.md`
- relevant release notes and integration summary artifacts

# Inputs

- accepted release decision
- released version or immutable artifact digest
- exposure target and completion boundary
- publication or handoff evidence when available
- release notes and integration summary references when applicable

# Must Do

- confirm exposure starts from an accepted release
- name the exposure extension or repository-specific exposure path
- record the exposure boundary reached
- record publication, handoff, deployment-manifest, and digest evidence without conflating their responsibilities
- state when production rollout is not yet proven
- route runtime, stakeholder, operator, or downstream signals into evidence for a later `extract` pass

# Must Not Do

- do not redefine the MRL core loop
- do not treat exposure as release acceptance
- do not claim production deployment without infrastructure-owned manifest or environment evidence plus immutable digest
- do not patch behavior directly from exposure feedback without first recording evidence and choosing the next loop entry point
- do not act as long-term operations ownership

# Outputs

- `work/changes/<id>/exposure.md` when exposure is planned, partial, completed, or blocked
- evidence references for publication, handoff, deployment truth, and follow-up extraction
