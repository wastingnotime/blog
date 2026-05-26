---
name: local-vault-access
description: Use when Codex needs a secret from the local KeePassXC vault before using gh, aws, or similar tools; reads CODEX_KC_VAULT and CODEX_KC_SECRET at runtime and never persists secret values.
---

# Local Vault Access

This repository owns the local-vault-access policy so secret handling stays consistent during repository work.

Use this skill when a task needs a secret that lives in the local KeePassXC vault.

## Inputs

- `CODEX_KC_VAULT`: path to the `.kdbx` database
- `CODEX_KC_SECRET`: vault password
- `keepassxc-cli`: command-line client for vault access

## Workflow

1. If the entry name is unknown, search first with the repo-local vault helper if it exists, or the equivalent local vault search path in the current environment.
2. Read the needed value with the repo-local vault helper when available, or the equivalent local vault read path in the current environment.
3. Prefer transient environment injection for commands that need the secret.
4. Keep the secret out of chat, logs, tracked files, and commit messages.
5. If unlock fails, verify the two env vars and confirm `keepassxc-cli` is installed.

## Usage patterns

- GitHub token:

  Use transient injection to validate `gh auth status` without persisting the token.

- AWS access key material:

  Use transient injection to run `aws sts get-caller-identity` without writing credentials to disk.

- Search the vault:

  Search for the entry name before assuming it is absent.

## Rules

- Do not print secret values unless the task explicitly needs them in stdout.
- Do not write decrypted secrets to repo files or generated reports.
- Do not ask the user to paste a secret that is already in the vault unless unlock or search fails.

