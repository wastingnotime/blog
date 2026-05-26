# Release Notes

- Change: `<change_id>`
- Date: YYYY-MM-DD
- Repository: `<repository>`
- Released version: `<version, commit sha, package version, or image tag>`
- Artifact digest: `<image digest, package digest, or n/a>`
- Status: draft | accepted | superseded

## Summary

Describe the delivered behavior and operational impact in concrete terms.

## Features

- `<delivered feature or capability>`

## Fixes

- `<delivered fix>`

## Infrastructure

- `<infra, packaging, deploy, runtime, or config change>`

## Contracts And Compatibility

- Breaking changes: `<yes/no and details>`
- Migrations: `<migration notes or n/a>`
- API changes: `<contract reference or n/a>`
- Event changes: `<contract reference or n/a>`

## Observability

- `<metrics, logs, traces, alerts, dashboards, or n/a>`

## Validation

- Integration summary: `<path, url, id, or pending>`
- Test evidence: `<path or summary>`

## Related

- Campaign: `<campaign id, url, or n/a>`
- Resolved blockers: `<blocker ids or n/a>`
- Findings addressed: `<finding ids or n/a>`

## Production Trace

- Infrastructure PR: `<url or pending>`
- Deployment manifest: `<path or pending>`
- Production image digest: `<digest or pending>`

Do not mark this release as production-deployed unless the production manifest state confirms the digest.
