# BTRBK Manage

## ⚠️ Instability Notice

**WARNING: This software is pre-v1.0.0 and should not be used without understanding the code first.**

Not all backup configs are supported yet. (E.g `snapshot_create on_change`). Use at your own risk and always test thoroughly before using on production data.

## Capabilities

- Delete specific snapshots/backups (later called "snaps")
- Mark snaps as "protected" to avoid deletion by automatic cleanup
- Restore snapshot

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
