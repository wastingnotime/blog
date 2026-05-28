# Repository Guidelines

## Project Structure & Module Organization
This repository is an MRL starter. Strategic docs live at the root and working design material lives under `docs/`. Read `docs/operating/mrl_reference.md` and `docs/operating/skills_workflow.md` before changing the workflow.

For cross-repository findings discovered locally, follow `docs/operating/cross_repo_findings_guidance.md`. Create findings in this repository, escalate multi-repository impact to management, and do not implement fixes owned by another repository.
For MRL-managed work resolved by code or configuration changes, follow `docs/operating/issue_pr_linking_guidance.md`. Open a pull request in this repository and link it back to the issue, campaign record, or finding record.

Use this structure as the default shape:

```text
src/app/{domain,application,infrastructure,interfaces}
tests/{unit,integration,builders,fixtures}
docs/{operating,building,evaluation,semantics,slices}
.agents/skills/
```

Record structural deviations in `decisions.md`.

## Build, Test, and Development Commands
Keep tooling lightweight until the first slice exists.

- `uv venv` creates the local virtual environment.
- `uv pip install -e .` installs the editable package.
- `uv run pytest` runs the test suite.
- `uv run python -m src.app.interfaces.cli.run_scenario` is the preferred shape for a local scenario runner when one exists.

## Coding Style & Naming Conventions
Prefer Python 3.12+, explicit types, 4-space indentation, and business-oriented names. Use verb-driven use cases such as `PlaceOrder` and intention-revealing repositories such as `get_by_id` and `save`.

## Testing Guidelines
Use tests as specification. Start with domain tests, add integration tests for mappings and end-to-end flows, and keep time, IDs, and external responses deterministic.

## Commit & Pull Request Guidelines
Use Conventional Commits for commit subjects, choosing an appropriate type such as `feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, or `chore`. Keep commits scoped to one slice or doc change, and include test evidence in pull requests.
After each completed change, create the commit before moving on to the next task.

<!-- mrl-cli agent hint start: wnt-overlay-hints -->
## WNT Overlay Hints

For managed-repository handoffs, follow `docs/operating/extensions/wnt/cross_repo_findings_guidance.md`.

For issue and pull request linking in MRL-managed work, follow `docs/operating/extensions/wnt/issue_pr_linking_guidance.md`.

For delivery or production-status questions, read `docs/operating/extensions/wnt/release_delivery_validation.md`.

This hint block is installed by `mrl-cli` from `mrl-extension-wnt`. It is not repository identity; root strategic docs and the rest of `AGENTS.md` remain owned by this repository.
<!-- mrl-cli agent hint end: wnt-overlay-hints -->