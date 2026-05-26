# Request Slice Map

- Change: `<change_id>`
- Request: `work/changes/<id>/request.md`
- Status: draft | ready-for-build | revised | accepted

## Request Boundary

Summarize the request intent that EGD and release should evaluate.

## Slice Mapping

| Slice | Status | Request coverage | Acceptance evidence |
| --- | --- | --- | --- |
| `docs/slices/<slice>.md` | planned | `<what part of the request this slice covers>` | `<tests, scenario evidence, contract artifact, or review point>` |

## Out Of Scope

- `<explicitly excluded request-adjacent behavior>`

## Open Questions

- `<question that must return to extract/refine before build or release>`

## EGD Notes

Describe what request-level expectation-gap detection must check across the listed slices.
