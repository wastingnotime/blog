---
name: campaign-management-protocol
description: Coordinate campaigns that span multiple repositories, regardless of where they start. Use when a repo change needs a source campaign issue, shared campaign memory, target-repository issues, Discord notifications, handoffs, status, decisions, and a final summary, including blocker-originated follow-up.
---

# Campaign Management Protocol

This repository owns the campaign-management protocol for cross-repository work.
See [`docs/operating/extensions/wnt/inter_repository_coordination_protocol.md`](../../docs/operating/extensions/wnt/inter_repository_coordination_protocol.md) for the shared channel model.

## When To Use

Use this skill when a change touches more than one repository or when a repository needs to start a cross-repo campaign.

Typical triggers:

- contract change that affects dependent repositories
- integration or runtime change with downstream impact
- finding that needs multi-repository follow-up
- Codex-triggered campaign kickoff from any managed repository
- blocker that produced one or more target issues in other repositories

## Workflow

1. Create or update the campaign issue in the originating repository. This is the source of truth for campaign state and completion.
2. Mirror the campaign into `management` as shared memory when that management repo exists in the workflow.
3. Use the canonical campaign issue body shape from [`references/campaign_issue_template.md`](../../references/campaign_issue_template.md).
4. Read or write the campaign intent and impact map.
5. Work one repository at a time unless explicitly told otherwise.
6. Create or update GitHub issues in each affected repository.
7. When the campaign comes from a blocker, ensure the local discovering issue and the target issue(s) are linked in both directions.
8. If a campaign task is resolved with code or configuration changes in an MRL-managed repository, open a pull request in the affected repository and link it back to the campaign issue and repository-local issue.
9. When the campaign issue is created, verify that the repository has GitHub Actions workflows for issue-open and issue-close Discord notifications.
10. If the workflows are missing, stop and ask the user to install them, pointing to the reference implementation at `infra-platform/.github/workflows/notify-discord-on-issue-open.yml` and `infra-platform/.github/workflows/notify-discord-on-issue-close.yml`, or to the repository's own notification guidance.
11. If the workflows exist, use them as the campaign notification path and keep campaign notices aligned with the campaign issue template and management mirror.
12. Post coordination notices in `#coordination` (`1502714561923383296`) when the campaign changes materially.
13. Close the campaign issue only after all linked repository issues are closed.
14. Ensure the issue-close notification path exists so campaign completion can be announced consistently.
15. Update campaign status and decisions as scope or ownership changes.
16. Finish with a concise final summary that records outcome, repository results, and any residuals.

## Handoff Rules

- Every affected repository gets a repository-specific handoff file.
- Keep each handoff bounded to that repository's scope and decision space.
- Do not widen campaign scope to unrelated fixes.
- If the campaign depends on a local finding, record the finding as evidence and coordinate through management instead of improvising a cross-repo fix.

## Management Boundary

Campaigns may start anywhere, but the campaign issue remains the durable source of truth for state and completion.

Codex creates or updates campaign issues; humans do not file campaign issues manually.

`management` mirrors that issue as the durable memory for campaign intent, status, decisions, and final summary when a management repo is part of the workflow.

When a campaign starts from a finding, use the `cross-repo-findings` skill and the supporting guidance in [`docs/operating/extensions/wnt/cross_repo_findings_guidance.md`](../../docs/operating/extensions/wnt/cross_repo_findings_guidance.md). The finding workflow stays authoritative for discovery and ownership; this skill only coordinates the campaign that follows.

## Discord Threading

When the Discord tooling supports it, keep campaign updates anchored to the original campaign notice by replying to or threading from that message instead of posting disconnected updates.

Use the same anchor for:

- material scope changes
- ownership changes
- repository handoff updates
- final completion notes
- record the anchor in the campaign issue and handoff notes so it can be reused for every follow-up
- the issue-open notification creates the anchor
- the issue-close notification replies to the anchor with the campaign completion or closure status

If replies or threads are not available, include the original message link in each follow-up so the campaign conversation remains connected.
The campaign packet should show the `Discord Channel` in its Discord Anchor section.

Target repository PR handling is defined in [`references/target_repository_pr_handling.md`](../../references/target_repository_pr_handling.md). Non-code resolution handling is defined in [`references/target_repository_resolution.md`](../../references/target_repository_resolution.md). False-positive handling is defined in [`references/target_repository_false_positive.md`](../../references/target_repository_false_positive.md).
