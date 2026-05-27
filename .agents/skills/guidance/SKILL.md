---
name: guidance
description: Help a repository owner or operator understand how to use MRL in this repository. Use when the owner has questions about phases, artifacts, skills, request/slice/release boundaries, packs, or workflow choices and needs explanatory guidance rather than repository edits.
---

# Mission

Answer MRL operating questions using the repository artifacts as the source of truth.

# Read First

- `docs/operating/mrl_reference.md`
- `docs/operating/skills_workflow.md`
- `docs/operating/mrl_starter.md`
- `docs/operating/best_practices.md`
- `docs/operating/packs.md`
- `decisions.md`
- relevant skill files under `.agents/skills/`
- relevant project artifacts named in the question

# Inputs

- owner or operator questions about MRL
- repository artifacts relevant to the question
- optional current change/request context

# Must Do

- answer from explicit repository artifacts and decisions
- explain which MRL phase or skill is appropriate for the situation
- distinguish current repository policy from general MRL advice
- point to artifacts that should be read or updated next
- identify uncertainty or missing decisions clearly
- keep advice practical and bounded to the question

# Must Not Do

- do not modify repository files
- do not run build, release, expose, or other phase work
- do not invent policy when `decisions.md` or operating docs are silent
- do not provide legal, financial, or compliance advice beyond pointing out repository evidence gaps
- do not rely on prior conversation as authoritative memory

# Outputs

- concise guidance answer
- recommended next artifact, phase, or decision when useful
