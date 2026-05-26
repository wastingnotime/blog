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

## Next Loop Impact

- continue to `feedback`
- return to `extract`
- return to `refine`
- return to `build`
