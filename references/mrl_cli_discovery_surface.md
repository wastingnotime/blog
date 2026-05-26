# MRL CLI Discovery Surface

## Purpose

This document describes the data contract that makes `mrl-extension-wnt` discoverable to `mrl-cli`.

The CLI owns command names, execution flow, and update policy. This repository owns the stable surface that the CLI reads.

## Discoverable Fields

- overlay name
- overlay source version
- overlay source reference
- installed overlay version
- lock file path
- update mode

## Sources

`mrl-cli` should be able to discover the overlay state from:

- `manifest.yaml`
- `mrl-extension.lock.yaml` in a consuming repository when present
- repository guidance that explains any intentional local divergence

## Expected Interpretation

- `manifest.yaml` is the source version for the overlay repository
- `mrl-extension.lock.yaml` records the installed overlay version in a consuming repository
- `provenance.lockfile` identifies the expected lock file name
- `provenance.update_mode` identifies whether the consuming repository expects pinned or tracked behavior
- `provenance.record_source_ref` tells `mrl-cli` whether it should preserve source revision information
