# WNT Pack Extension: terraform_aws_environment

## Purpose

Use this pack when the primary work surface is Terraform-managed AWS infrastructure.

## Scope

This pack extension defines WNT-owned implementation defaults.
It does not define MRL core behavior or project-specific domain behavior.

This pack exists to define reusable implementation defaults for repositories whose main behavior is infrastructure definition, migration, and controlled environment evolution rather than application runtime logic.

---

## Shape

- language: HCL with optional helper scripts
- runtime topology: Terraform root modules plus reusable modules
- target platform: AWS
- operating mode: migration first or controlled evolution of an existing environment

---

## Good Fit

- AWS environments managed primarily through Terraform
- repositories that need explicit migration from a legacy Terraform codebase
- infrastructure work where plans, validation, and parity matter more than application-layer abstractions

---

## Less Suitable

- application repositories where Terraform is incidental
- multi-cloud platforms where AWS is only one small concern
- repositories whose primary behavior is a long-running service rather than infrastructure definition

---

## Layout

```text
terraform/
  environments/
    <environment>/
  modules/
    <module_name>/

tests/
  unit/
  integration/
  fixtures/

docs/
  semantics/
  slices/
```

---

## Boundaries

- infrastructure intent lives in `docs/semantics/` and `docs/slices/`
- Terraform root modules live under `terraform/environments/`
- reusable infrastructure patterns live under `terraform/modules/`
- helper scripts should stay small and operationally narrow
- durable environment truth should live in Terraform and tracked repository artifacts, not only in shell history

---

## Testing Defaults

- unit/spec focus: helper scripts and generated configuration behavior when present
- integration focus: module wiring and root-module validation
- determinism rules: pin provider behavior through lock files where appropriate and avoid hidden manual state changes
- preferred test command: `terraform validate`

---

## Scenario Runner Shape

- preferred runner location: `Makefile` or small scripts under `scripts/`
- preferred command: `terraform plan`
- runner responsibility: validate one bounded infrastructure slice with explicit variables and environment scope
- runner must not: hide destructive changes or blur environment boundaries

---

## Runtime Targeting Rules

- slices using this pack should declare runtime targets as:
  - `terraform/<environment>/<root>`
  - `aws/<service>` when a slice is tied to a specific platform surface
- cross-root dependencies should be made explicit in slice docs and module boundaries
- operational contracts should be expressed in tracked docs or configuration, not left implicit

---

## Naming Conventions

- root modules: intention-revealing environment or capability names
- reusable modules: capability-oriented names
- variables and outputs: explicit infrastructure terms
- helper scripts: verb-first operational names

---

## Rules

- extract existing environment facts before redesigning resources
- move one bounded area at a time
- validate with `terraform fmt`, `terraform validate`, and targeted `terraform plan`
- keep authority boundaries and accepted drift explicit in tracked artifacts

---

## Starter Tooling

- package manager: n/a
- formatter/linter: `terraform fmt`
- test framework: `terraform validate` plus targeted plan review
- persistence default: Terraform state backend selected by the adopting repository

---

## Notes

This pack is a reusable default, not a law. If a repository diverges meaningfully, record that in `decisions.md`.
