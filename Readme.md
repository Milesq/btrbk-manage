# BTRBK Manage

A unified management tool for organizing, protecting, and cleaning up btrbk snapshots.

## ⚠️ Instability Notice

**WARNING: This software is pre-v1.0.0 and should not be used without understanding the code first.**

Not all backup configs are supported yet. (E.g `snapshot_create on_change`). Use at your own risk and always test thoroughly before using on production data.

## Architecture

All functionality is consolidated into a single `btrbk-manage` binary that provides an interactive TUI for managing snapshots.

## Usage

```bash
btrbk-manage [snapshot-directory]
```

The tool launches an interactive terminal UI where you can:
- Navigate snapshots/backups using arrow keys
- Press `Space` to protect/unprotect backups
- Press `Enter` to edit protection notes
- Press `T` to review and purge the trash
- Press `q` to quit

## Capabilities

- Delete specific snapshots/backups (later called "snaps")
- Mark snaps as "protected" to avoid deletion by automatic cleanup
- Restore snapshots
- Review and purge trashed backups

## Terminology

This project uses specific terminology that differs from standard btrbk:

- **Snapshot**: An individual btrfs snapshot (same as btrbk)
- **Backup**: A logical group of snapshots taken at the same point in time across multiple volumes. For example, all snapshots created at `20250909T2343` for volumes `@`, `@home`, and `@srv` form one "backup". This is NOT the same as btrbk's remote backup feature.

<!-- Note: Not all snapshots in a backup group are necessarily created at the exact same time. If a snapshot is missing for a volume at the backup timestamp, the group will include the most recent previous snapshot of that volume instead. -->

## Configuration

The tool reads snapshot directories and utilizes `btrbk list` to enumerate available snapshots.

## Features

### Deleting Snaps

- By default, snaps are moved to `SNAPDIR/.trash` instead of being permanently deleted
- Trashed items can be reviewed and purged from within the UI (press `T`)
- Provides a safety net for accidental deletions

### Protecting Snaps

The tool creates the following metadata structure in the snapshots directory:
```.meta
└── 20250909T2343
    ├── info.yaml
    └── snaps
        ├── @
        ├── @home
        └── @srv
```

When unchecking a protected backup, it is moved to the `.trash` directory instead of being immediately deleted.

You can review and purge the trash by pressing `T` within the UI.

`info.yaml` contains:
- User notes for the backup
- Restoration dates
- Tags

### Unified Snaps

Manager may group snaps by date.

if `snapshot_create onchange` option is applied,
