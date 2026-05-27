---
name: adoption-diagnose
description: Diagnose whether an adopting repository has intentionally replaced starter defaults with project-specific artifacts. Consider on the first pass through repository guidance before substantial project-specific work, and use when the repository owner wants a readiness check for licensing, README content, selected pack, starter placeholders, semantic artifacts, or other adoption residue without changing production code.
---

# Mission

Audit whether the repository has moved from MRL starter state into an intentionally adopted project state.

# Read First

- `readme.md`
- `LICENSE`
- `decisions.md`
- `architecture.md`
- `groundrules.md`
- `docs/operating/mrl_starter.md`
- `docs/operating/skills_workflow.md`
- `docs/building/project_structure.md`
- `docs/semantics/model_hypothesis.md`
- `docs/semantics/domain_background_knowledge.md`
- relevant `docs/slices/` and `work/changes/` artifacts when work has started

# Inputs

- current repository artifacts
- repository owner questions about adoption readiness
- optional target licensing or project-positioning intent from the owner

# Must Do

- check whether licensing has an explicit accepted decision in `decisions.md`
- check whether `LICENSE` still reflects the intended repository licensing approach
- check whether `readme.md` describes the adopting project rather than only MRL itself
- check whether root strategic docs are project-specific enough for the current adoption stage
- check whether semantic docs and slice/change artifacts still contain starter placeholders
- check whether the selected implementation pack is recorded or intentionally replaced
- classify observations by severity: blocker, warning, or note
- write or update `work/adoption_diagnosis.md` when a durable report is requested or when observations are non-trivial
- recommend the next MRL phase or owner decision when adoption gaps are found

# Must Not Do

- do not change licensing text, README content, or project strategy directly unless the owner explicitly asks for edits
- do not provide legal advice; identify repository evidence and decision gaps only
- do not write production code
- do not run the canonical MRL loop phases on behalf of this audit
- do not rely on conversational memory over repository artifacts

# Outputs

- adoption readiness observations
- optional `work/adoption_diagnosis.md`
- recommended follow-up decisions or MRL phases
