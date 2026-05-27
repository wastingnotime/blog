# WNT Issue / PR Linking Guidance

## Purpose

Define the local resolution shape when a tracked issue, campaign task, blocker, or finding is answered by repository changes in a WNT-managed repository.

## Scope

This document is WNT overlay guidance.
It applies to repositories managed under the WNT MRL contract.

It does not require a PR for outcomes that are resolved entirely by documentation, discussion, triage, or decision without repository changes.

Use it with:

- `docs/operating/extensions/wnt/inter_repository_coordination_protocol.md`
- `docs/operating/extensions/wnt/cross_repo_findings_guidance.md`
- `.agents/references/target_repository_pr_handling.md`

---

## Core Rule

When repository changes are required:

- open a pull request in the repository that owns the change
- link the originating issue, campaign issue, blocker, or finding issue in the PR body
- link the PR back from the originating issue, campaign record, blocker, or finding record
- keep the issue or record as the tracking artifact
- keep the PR as the implementation artifact

Prefer explicit bidirectional linkage over relying on implied closure alone.

If the repository uses GitHub closing keywords, they are additive. They do not replace the explicit link from the issue or record back to the PR.

---

## Authority

This document owns the WNT issue-to-PR linkage rule.

`docs/operating/extensions/wnt/inter_repository_coordination_protocol.md` owns channels, urgency, and coordination routing.
`docs/operating/extensions/wnt/cross_repo_findings_guidance.md` owns findings and blocker workflow.

If this document and a workflow-specific document disagree about when to create a target issue, the workflow-specific document wins.
If this document and the coordination protocol disagree about channel or urgency, the coordination protocol wins.

---

## Target Repository Handling

When a target repository issue requires code or configuration changes:

- assign the issue promptly according to the repository's ownership rules
- open the PR in the target repository
- include the target issue link in the PR body
- include the closing keyword in the PR body when automatic closure is desired
- comment on the issue with the PR link when the PR is opened
- do not close the issue before the implementation conclusion is recorded

The target issue remains the durable statement of ownership and resolution.
The PR remains the implementation review and merge artifact.

---

## Campaign And Finding Records

When the originating work is a campaign, blocker, or finding:

- keep the campaign, blocker, or finding record linked to every target issue and PR it caused
- keep the target issue linked back to the origin record
- close the origin only after linked target issues are closed or explicitly rejected with evidence
- preserve rejected ownership decisions instead of deleting the linkage trail

This keeps coordination memory separate from implementation review while making the chain easy to follow.

---

## Notes

- A PR without a linked issue is acceptable only for small local maintenance that has no originating tracked record.
- A closing keyword is useful automation, not the source of truth.
- For cross-repository findings, follow `docs/operating/extensions/wnt/cross_repo_findings_guidance.md` first and use this document for the issue-to-PR linkage shape.
