---
name: create-managed-repository
description: Create a new Wasting No Time GitHub repository from `wastingnotime/mrl-starter`, install `mrl-extension-wnt`, create the local lets launcher, and announce the repository to Discord `#management` without writing to the management repository. Use when an MRL-enabled WNT repository is asked to create, bootstrap, clone, or introduce a new `wastingnotime` MRL repository.
---

# WNT Managed Repository Creation

## Overview

Use this workflow from any WNT repository that has the `mrl-extension-wnt`
overlay installed.

The workflow creates a new `wastingnotime` repository from
`wastingnotime/mrl-starter`, makes the local clone usable, installs the WNT
overlay with `mrl-cli`, creates the matching lets launcher, commits and pushes
the target bootstrap, then sends a repository-introduction message to Discord
`#management` (`1509567107891986644`).

Do not write to the `management` repository. Management owns its own
classification and recording decisions after it receives the introduction
message.

## Required Context

Start from the repository where the user requested creation. Read its local
guidance before acting:

```bash
sed -n '1,220p' AGENTS.md
```

Use the installed WNT coordination model when present:

```bash
sed -n '1,220p' docs/operating/extensions/wnt/inter_repository_coordination_protocol.md
```

Use `mrl-cli` when it is available. If it is not available, stop before
creating the repository unless the user explicitly approves manual bootstrap
steps.

## Inputs

Infer or confirm:

- repository name, for example `landing-page`
- visibility, usually `--private` unless the user asks for public
- purpose, for the GitHub description and Discord introduction
- overlay source path, defaulting to
  `/home/henrique/repos/github/wastingnotime/mrl-extension-wnt`
- local clone path, defaulting to
  `/home/henrique/repos/github/wastingnotime/<repo>`

Honor explicit visibility requests. If the requested visibility seems
surprising for the stated purpose, warn briefly and continue only when the
user confirms or the request is already explicit.

## Create Repository

Check current state:

```bash
git status --short
gh auth status
gh repo view wastingnotime/<repo> --json nameWithOwner,isPrivate,url,defaultBranchRef
gh repo view wastingnotime/mrl-starter --json nameWithOwner,isTemplate,isPrivate,url,defaultBranchRef
test -d /home/henrique/repos/github/wastingnotime/<repo> && git -C /home/henrique/repos/github/wastingnotime/<repo> status --short || true
mrl-cli --help
```

Create the repository from the template when the target does not exist:

```bash
gh repo create wastingnotime/<repo> \
  --private \
  --template wastingnotime/mrl-starter \
  --description "<short purpose>" \
  --clone=false
```

Use `--public` instead of `--private` only when intended.

Verify:

```bash
gh repo view wastingnotime/<repo> --json nameWithOwner,isPrivate,isTemplate,url,defaultBranchRef,description,templateRepository
```

## Clone And Prepare

Clone to the standard path:

```bash
git clone git@github.com:wastingnotime/<repo>.git /home/henrique/repos/github/wastingnotime/<repo>
```

Inspect:

```bash
git -C /home/henrique/repos/github/wastingnotime/<repo> status --short --branch
git -C /home/henrique/repos/github/wastingnotime/<repo> log --oneline -5
sed -n '1,200p' /home/henrique/repos/github/wastingnotime/<repo>/AGENTS.md
```

If the clone lacks author identity, copy the standard repo identity used for
new WNT repositories:

```bash
git -C /home/henrique/repos/github/wastingnotime/<repo> config user.name "hriccio"
git -C /home/henrique/repos/github/wastingnotime/<repo> config user.email "hriccio@wastingnotime.org"
```

## Install WNT Overlay

Install `mrl-extension-wnt` after the template repository is cloned:

```bash
mrl-cli install \
  --source /home/henrique/repos/github/wastingnotime/mrl-extension-wnt \
  --target /home/henrique/repos/github/wastingnotime/<repo>
```

Verify the installed state:

```bash
mrl-cli discover \
  --source /home/henrique/repos/github/wastingnotime/mrl-extension-wnt \
  --target /home/henrique/repos/github/wastingnotime/<repo>
git -C /home/henrique/repos/github/wastingnotime/<repo> status --short
```

## Create Lets Launcher

Create a matching launcher at:

```text
/home/henrique/scripts/lets-<repo>
```

Use the existing launcher convention:

```sh
#!/bin/sh

. "$HOME/scripts/lets-common.sh"
launch_lets "lets-<repo>" "$HOME/repos/github/wastingnotime/<repo>"
```

Make it executable:

```bash
chmod +x /home/henrique/scripts/lets-<repo>
```

Verify against nearby launchers when unsure:

```bash
sed -n '1,80p' /home/henrique/scripts/lets-mrl-starter
sed -n '1,80p' /home/henrique/scripts/lets-runtime-sandbox
ls -l /home/henrique/scripts/lets-<repo>
```

## Commit Target Bootstrap

Commit and push the target repository bootstrap after the overlay install and
launcher verification:

```bash
git -C /home/henrique/repos/github/wastingnotime/<repo> status --short
git -C /home/henrique/repos/github/wastingnotime/<repo> add .
git -C /home/henrique/repos/github/wastingnotime/<repo> commit -m "chore: bootstrap WNT MRL repository"
git -C /home/henrique/repos/github/wastingnotime/<repo> push origin main
```

Record the bootstrap commit SHA:

```bash
git -C /home/henrique/repos/github/wastingnotime/<repo> rev-parse --short HEAD
```

If the template and overlay produce no target repository changes, record the
current `HEAD` SHA and explain that no bootstrap commit was needed.

## Announce To Management

Send a Discord message to `#management` (`1509567107891986644`) after the
target repository is pushed. Use the available Discord sender in the local
environment, such as `discordctl.py`, an approved Discord MCP tool, or the
current WNT notification mechanism.

The message should introduce the new repository without assigning management
an action:

```text
New WNT MRL repository created

Repository: wastingnotime/<repo>
URL: https://github.com/wastingnotime/<repo>
Visibility: private/public
Purpose: <user-provided purpose>
Template source: wastingnotime/mrl-starter
WNT overlay: mrl-extension-wnt <source ref if available>
Local path: /home/henrique/repos/github/wastingnotime/<repo>
Lets launcher: /home/henrique/scripts/lets-<repo>
Bootstrap commit: <sha or no-new-commit note>
Created from: <origin repository path or name>
```

Do not include a requested management action. The management repository decides
what, if anything, to do with the information.

## Verify

Run:

```bash
git -C /home/henrique/repos/github/wastingnotime/<repo> status --short --branch
git -C /home/henrique/repos/github/wastingnotime/<repo> ls-remote origin main
gh repo view wastingnotime/<repo> --json nameWithOwner,isPrivate,url,defaultBranchRef,description,templateRepository
mrl-cli discover --source /home/henrique/repos/github/wastingnotime/mrl-extension-wnt --target /home/henrique/repos/github/wastingnotime/<repo>
ls -l /home/henrique/scripts/lets-<repo>
```

The source repository running this workflow should remain clean unless the
user explicitly asked for local changes there.

## Final Response

Report:

- GitHub repository URL
- visibility
- template source
- WNT overlay source
- local clone path
- lets launcher path
- target repository bootstrap commit or no-new-commit note
- Discord `#management` message status
- any unpushed state in the new repository
