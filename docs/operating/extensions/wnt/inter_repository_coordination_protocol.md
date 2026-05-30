# Inter-Repository Coordination Protocol

## Purpose

This repository uses a four-channel coordination model for the WNT ecosystem.

The goal is to separate:

- slow coordination
- active dependency interruption
- runtime incident response
- management visibility

## Authority

This document owns the WNT channel model and coordination routing rules.

Workflow-specific documents, such as `docs/operating/extensions/wnt/cross_repo_findings_guidance.md`, may define how to execute a specific workflow, but they must defer to this protocol for channel selection, urgency, and coordination ownership.

## Channels

| Channel | ID | Purpose |
| --- | --- | --- |
| `#coordination` | `1502714561923383296` | Cross-repository coordination, campaigns, findings, migrations, compatibility notices, and architecture evolution |
| `#blockers` | `1502715067186024628` | Active dependency interruption between repositories |
| `#ops-alerts` | infra-platform contract | Runtime infrastructure and production remediation |
| `#management` | `1509567107891986644` | Repository introductions and management-visible facts that management may classify or record at its own boundary |

## Meaning

### `#coordination`

Use for:

- campaigns
- findings
- migrations
- compatibility notices
- architecture evolution
- cross-repository decisions
- integration observations
- planned changes
- dependency awareness
- contract evolution notices

This channel is asynchronous, context-rich, durable, and low urgency.

### `#blockers`

Use for:

- blocked repository execution
- unresolved dependency contracts
- incompatible integrations preventing progress
- CI/CD progression blockers
- missing required changes in another repository
- urgent dependency resolution requests

This channel is operational, actionable, high signal, and interruption oriented.

### `#ops-alerts`

Use for:

- production incidents
- swarm degradation
- deployment failures
- telemetry alerts
- infrastructure instability
- runtime regressions
- operational remediation workflows

This channel is runtime focused and interruption critical.

The `#ops-alerts` contract is owned by `infra-platform` and is delivered through Grafana alerts and Air-bot. This repository should treat it as an external runtime-alerting integration and should not open coordination messages there directly.

### `#management`

Use for:

- introducing newly created WNT MRL repositories
- sharing management-visible repository facts
- surfacing context that management may classify, record, ignore, or follow up on at its own boundary

This channel is informational. Messages sent here must not instruct management
to classify, record, or take a specific action unless a separate management
workflow explicitly requested that action.

## Relationship To Repository Issues

- findings belong in `#coordination` (`1502714561923383296`)
- campaigns belong in `#coordination` (`1502714561923383296`)
- blockers belong in `#blockers` (`1502715067186024628`)
- runtime incidents belong in `#ops-alerts` through the infra-platform alerting contract, not through local repo workflows
- new repository introductions belong in `#management` (`1509567107891986644`)

## Discord Anchors

When the tooling supports it, create one anchor message per issue or campaign and reply to that anchor for follow-up updates.

The issue or campaign packet should record the Discord channel in its Discord Anchor section:

- `#coordination` (`1502714561923383296`) for findings and campaigns
- `#blockers` (`1502715067186024628`) for blockers
