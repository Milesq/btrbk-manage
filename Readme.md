# BTRBK Manage

A management tool for organizing, protecting, and cleaning up btrbk snapshots.

## ⚠️ Instability Notice

**WARNING: This software is pre-v1.0.0 and should not be used without understanding the code first.**

Not all backup configs are supported yet. (E.g `snapshot_create on_change`). Use at your own risk and always test thoroughly before using on production data.

## Capabilities

- Delete specific snapshots/backups (later called "snaps")
- Mark snaps as "protected" to avoid deletion by automatic cleanup
- Restore snapshot

## Terminology

This project uses specific terminology that differs from standard btrbk:

- **Snapshot**: An individual btrfs snapshot (same as btrbk)
- **Backup**: A logical group of snapshots taken at the same point in time across multiple volumes. For example, all snapshots created at `20250909T2343` for volumes `@`, `@home`, and `@srv` form one "backup". This is NOT the same as btrbk's remote backup feature.

<!-- Note: Not all snapshots in a backup group are necessarily created at the exact same time. If a snapshot is missing for a volume at the backup timestamp, the group will include the most recent previous snapshot of that volume instead. -->

## General Config

- takes a snapshot directory utilizing `btrbk list`

## Deleting Snaps

- by default trashes snaps to SNAPDIR/.trash
- can permanently delete with `--permanent` flag
- [Unified Snaps] are disabled by default so you delete only selected snaps - enable them with --unified
-

## Protecting Snapsw

Manager creates following structure in the snapshots dir
```.meta
└── 20250909T2343
    ├── info.toml
    └── snaps
        ├── @
        ├── @home
        └── @srv
```

info.toml contains:
<!-- TODO:  -->
- additional user notes

### Unified Snaps

Manager may group snaps by date.

if `snapshot_create onchange` option is applied,
