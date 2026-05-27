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

<!-- mrl-cli patch start: AGENTS.md -->
# Repository Guidelines

## Scope
This file guides contributors working in this repository or an adopting project instance. It does not define MRL core behavior; core workflow guidance lives in `docs/operating/`.

## Project Structure & Module Organization
This repository is a private WNT extension overlay for MRL. Strategic docs live at the root and overlay material lives under `docs/`, `skills/`, `references/`, `templates/`, `agents/`, `contracts/`, `observability/`, `infra/`, and `codex/`.

Root strategic docs describe the current repository or adopting project instance. MRL core behavior lives in `docs/operating/` and should stay generic, portable, and operationally agnostic.

On the first pass through this repository's guidance, review `architecture.md`, `groundrules.md`, and the current overlay files before substantial project-specific work.
For delivery or production-status questions, read `docs/operating/extensions/wnt/release_delivery_validation.md` and treat infra-owned deployment manifests plus immutable artifact digests as production truth.

Use this structure as the default overlay shape:

```text
skills/
references/
templates/
agents/
contracts/
observability/
infra/
codex/
docs/operating/
```

Record structural deviations in `decisions.md`.

## Build, Test, and Development Commands
Keep tooling lightweight until the first overlay slice exists.

- `git status --short` inspects pending edits before commits.
- `git diff --check` catches whitespace and patch formatting issues.
- `mrl-cli --help` validates the installer/overlay CLI when it is available in the local environment.
- `mrl-cli sync --source ../mrl-extension-wnt --target .` refreshes installed WNT guidance in a consuming repository when run from that repository.

## Coding Style & Naming Conventions
Prefer Python 3.12+, explicit types, 4-space indentation, and business-oriented names. Use verb-driven use cases such as `PlaceOrder` and intention-revealing repositories such as `get_by_id` and `save`.

## Testing Guidelines
Use tests as specification. Start with domain tests, add integration tests for mappings and end-to-end flows, and keep time, IDs, and external responses deterministic.

## Commit & Pull Request Guidelines
Use Conventional Commits for commit subjects, choosing an appropriate type such as `feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, or `chore`. Commit after every completed and verified change before starting unrelated work. Keep commits scoped to one request, slice, or doc change, and include test evidence in pull requests.
<!-- mrl-cli patch end: AGENTS.md -->