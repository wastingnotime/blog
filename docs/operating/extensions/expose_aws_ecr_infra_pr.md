# Expose Extension: aws_ecr_infra_pr

## Purpose

Use this expose extension when a repository normally exposes accepted artifacts by publishing a container image through GitHub Actions to AWS ECR and then handing deployment intent to a separate infrastructure repository through a pull request.

Use `docs/operating/release_delivery_validation.md` with this extension when answering what is in production. The infrastructure repository deployment manifest plus immutable image digest is production truth; this extension records the application-side publication and handoff path.

## Scope

This document is an optional exposure overlay.
It does not define MRL core behavior or a mandatory project deployment path.

This extension keeps MRL core generic while making a recurring remote-infrastructure exposure path explicit and reusable.

---

## Good Fit

Use this extension when:

- the released artifact is a container image
- GitHub Actions is the normal automation surface
- AWS ECR is the image registry
- a separate infrastructure repository owns deployment manifests, params, or environment rollout
- repository responsibility ends at explicit handoff or at another clearly documented boundary

Do not use this extension when:

- the repository deploys directly without an infrastructure handoff repository
- the exposed artifact is not a container image
- the exposure path depends on a different deployment surface such as direct Terraform apply or a platform-native release controller

---

## Exposure Boundary

The default target is the remote environment represented by the infrastructure repository and its deployment configuration.

For repositories using this extension, exposure completion should be made explicit per repository. Typical boundaries are:

- opening the infrastructure PR
- updating an existing infrastructure PR
- observing that the infrastructure PR merged
- observing that the target environment rolled out the new artifact

The repository must state which of those counts as "done" in `decisions.md` and in each `exposure.md`.

Even when a repository chooses "opening the infrastructure PR" as its exposure boundary, that does not mean the change is in production. Production status requires the infra-platform manifest or environment state to show the immutable digest.

---

## Required Inputs

- `work/changes/<id>/release_decision.md`
- released version identifier
- container image naming contract
- Docker build context
- AWS account and region
- ECR repository target
- automation identity assumptions for GitHub Actions and AWS credentials
- infrastructure repository reference
- deploy parameter contract such as image URI, digest, tag, or release metadata fields
- release notes reference for the delivered artifact
- integration summary reference when ecosystem validation applies

---

## Outputs

- GitHub Actions workflow run URL or identifier
- published ECR image URI
- immutable image digest
- infrastructure repository branch or PR URL
- summary of changed deploy params
- release notes reference
- integration summary reference
- production deployment manifest reference when rollout is observed
- explicit exposure result:
  - handoff created
  - handoff updated
  - handoff blocked
  - downstream rollout observed
  - downstream rollout not yet observed

---

## Invariants

- exposure starts only from an accepted released version
- the image published to ECR corresponds to the accepted released version
- the infrastructure handoff updates explicit deploy inputs rather than relying on hidden mutable state
- the exposure artifact records the boundary between:
  - artifact publication
  - infrastructure handoff
  - actual environment rollout
- production claims are based on the infrastructure deployment manifest and immutable image digest
- the extension does not silently change domain or runtime semantics while preparing exposure

---

## Failure Classes

- GitHub Actions workflow failed before image publication
- Docker build succeeded locally but failed in CI
- ECR push failed or produced an unexpected target reference
- image tag and digest contract mismatch
- infrastructure repository parameter contract mismatch
- infrastructure PR could not be created or updated
- downstream rollout was intentionally deferred or blocked

---

## Evidence Sources

- GitHub Actions workflow run
- build logs and workflow artifacts
- ECR image URI and digest
- infrastructure repository PR diff
- merge status or deployment status from the infrastructure repository when available
- production deployment manifest path and image digest when rollout is observed
- release notes and integration summary references
- smoke or health evidence when the downstream environment exposes it

---

## Recommended Execution Shape

1. confirm that `release_decision.md` accepts the current version
2. build the runtime artifact in GitHub Actions
3. publish the artifact to ECR
4. capture the immutable image digest
5. record or update release notes for the delivered artifact
6. attach or reference the integration summary when ecosystem validation applies
7. update the infrastructure repository deploy params with explicit artifact references
8. open or update the infrastructure PR
9. record what boundary was reached and which evidence was observed
10. if rollout is observed, record the production deployment manifest reference and digest

---

## Repository Decision Template

```md
## DEC-XXXX - Default Exposure Extension

- Date: YYYY-MM-DD
- Status: accepted
- Owners: both

### Context
This repository usually exposes accepted artifacts through GitHub Actions, AWS ECR, and an infrastructure-repository handoff. That path should be reusable without coupling MRL core or the selected implementation pack to one platform workflow.

### Decision
This repository adopts `aws_ecr_infra_pr` as its default expose extension.

Released artifacts are exposed by:

- building a container image in GitHub Actions
- publishing that image to AWS ECR
- opening or updating an infrastructure repository PR with the explicit deploy parameters needed for rollout

Exposure completion for this repository means `<define repository-specific boundary here>`.

Production status is not inferred from release acceptance, image publication, or PR creation. It is confirmed only by the infrastructure deployment manifest or environment state that references the immutable image digest.

### Consequences
Exposure evidence becomes repeatable and inspectable across slices. The repository must keep the ECR publication contract and infrastructure handoff contract explicit as those systems evolve.

### Alternatives considered
Treat exposure as ad hoc per change, or fold this mechanism into an implementation pack. These were rejected because the concern is lifecycle-specific rather than implementation-specific.
```

---

## `exposure.md` Notes For This Extension

When this extension is active, each `work/changes/<id>/exposure.md` should record at least:

- extension name: `aws_ecr_infra_pr`
- workflow run URL or identifier
- ECR image URI
- image digest
- infrastructure PR URL
- exposure boundary reached
- gaps between publication, infra handoff, and actual rollout
- release notes reference
- integration summary reference
- production manifest and digest if rollout was observed
