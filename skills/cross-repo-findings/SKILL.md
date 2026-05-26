---
name: cross-repo-findings
description: Report a finding or blocker from any repository, create a local record, open repository issues, and notify the #coordination or #blockers channel. Use when a repository discovers a contract mismatch, integration failure, runtime inconsistency, blocker, or other issue that may affect another repository.
---

# Cross Repo Findings

This repository owns the cross-repo findings workflow. The durable repository guidance lives in [`docs/operating/cross_repo_findings_guidance.md`](../../docs/operating/cross_repo_findings_guidance.md).
See [`docs/operating/inter_repository_coordination_protocol.md`](../../docs/operating/inter_repository_coordination_protocol.md) for the shared channel model.

## When To Use

Use this skill when a repository discovers a problem that needs durable tracking, especially when the issue may impact another repository.

Typical triggers:

- contract mismatch
- integration failure
- runtime inconsistency
- deploy/runtime side effect
- unexpected behavior from another repository
- finding that is non-blocking but outside the current repository's ownership
- blocker that halts current work and likely belongs to another repository

## Workflow

1. Create a local finding record at `work/findings/<date>-<short-name>/finding.md`.
2. Capture observed behavior, expected behavior, impact, ownership, and evidence.
3. If the issue is a finding, make the non-blocking nature and the likely owning repository explicit.
4. If the issue is a blocker, make the blocked activity explicit and name the likely target repository even when ownership is not fully confirmed.
5. Open a GitHub issue in the discovering repository so the issue is visible in the repo that observed it.
6. When the issue is a finding or blocker and a likely owning repository is known, open a target issue there as well and link the local issue to the target issue.
7. Post a short notification to Discord `#coordination` (`1502714561923383296`) for findings or `#blockers` (`1502715067186024628`) for blockers.
8. Use the canonical issue body shape from [`references/finding_issue_template.md`](../../references/finding_issue_template.md) when the repository has that reference available.
9. If the issue affects more than one repository, escalate through management and create or update a campaign there.
10. If the issue is resolved with code or configuration changes in an MRL-managed repository, open a pull request in the affected repository and link it back to the discovering repository issue and local finding record.
11. Close the discovering repository issue when the issue is resolved or explicitly concluded.
12. Send a completion message to `#coordination` for findings or `#blockers` for blockers when that issue closes.

## Local Finding Record

Keep the finding precise and evidence-led.

Include:

- context
- observed behavior
- expected behavior
- impact
- issue type: finding | blocker
- suspected source
- evidence
- owning repository
- local impact

Do not prescribe implementation details for another repository.

## GitHub Issue

Create the issue in the discovering repository unless a repo-local workflow already exists that routes findings elsewhere.

Keep the issue short:

- title: short finding summary
- body: link to the local finding record, summarize evidence, identify the likely owning repository, and note any immediate local impact

For findings, the body should also state:

- that the issue is non-blocking
- that it is outside the current repository's ownership
- that it should be solved eventually by the owning repository

For blockers, the body should also state:

- the blocked activity
- whether the target repository is confirmed or only likely
- that the target issue is linked from the main discovering repository issue

The issue stays open until the finding is resolved or explicitly concluded.

Codex creates or updates discovering-repository finding issues; humans do not file finding issues manually.

## Discord Notification

Post a concise message to `#coordination` (`1502714561923383296`) for findings or `#blockers` (`1502715067186024628`) for blockers with:

- finding title
- discovering repository
- likely owning repository
- link or path to the local finding record
- Discord anchor message link or ID

Post a completion message to `#coordination` (`1502714561923383296`) for findings or `#blockers` (`1502715067186024628`) for blockers when the discovering repository issue closes.

When the tooling supports it, keep the Discord conversation anchored to the original finding or blocker message by replying in-thread or replying to the message instead of starting a disconnected post.
Record the anchor and Discord Channel in the local finding record and reuse them for every follow-up.
The issue-open notification should create the anchor.
The issue-close notification should reply to the anchor with the conclusion or resolution status.
The finding or blocker packet should show the Discord Channel in its Discord Anchor section.

Target repository PR handling is defined in [`references/target_repository_pr_handling.md`](../../references/target_repository_pr_handling.md). Non-code resolution handling is defined in [`references/target_repository_resolution.md`](../../references/target_repository_resolution.md). False-positive handling is defined in [`references/target_repository_false_positive.md`](../../references/target_repository_false_positive.md).

## Escalation

If the finding spans multiple repositories, use management campaigns for coordination after the local issue and Discord notification exist.

For findings, create the target issue when the owning repository is known or reasonably inferred, even if the current work is not blocked.

For blockers, create the likely target issue immediately when the target repository is reasonably clear, even if the owner has not yet confirmed responsibility.

Do not implement another repository's fix locally.

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

When the tooling supports it, reply in Discord to the original discovery message for ownership confirmations, revision requests, and completion notes.
