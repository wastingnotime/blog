# Integration Summary

- Change: `<change_id>`
- Date: YYYY-MM-DD
- Summary owner: `<repository or operator>`
- Result: passed | failed | partial | blocked
- Environment: `<integration-sandbox, local stack, staging, or equivalent>`

## Validated Matrix

| Repository | Version | Artifact digest |
|---|---|---|
| `<repository>` | `<sha, tag, or version>` | `<digest or n/a>` |

## Scenarios

| Scenario | Result | Evidence |
|---|---|---|
| `<scenario>` | passed | `<log, trace, screenshot, or artifact>` |

## Evidence

- Logs: `<url, path, or n/a>`
- Traces: `<url, path, or n/a>`
- Screenshots: `<url, path, or n/a>`
- Workflow run: `<url or n/a>`

## Findings And Blockers

- Findings opened: `<finding ids or n/a>`
- Blockers opened: `<blocker ids or n/a>`
- Campaign impact: `<campaign id or n/a>`

## Release Impact

State whether this validation supports release, blocks release, or requires another loop.

This summary validates a version. It is not a deployment manifest, release note, or campaign record.
