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
- Press `t` to review and purge the trash
    - D for deleting all trashed backups
- Press `m` to filter protected backups the trash
- Press `q` to quit

## Capabilities

- Delete specific backups
- Mark snaps as "protected" to avoid deletion by automatic cleanup
- Restore snapshots
- Review and purge trashed backups

## Terminology

This project uses specific terminology that differs from standard btrbk:

- **Snapshot**: An individual btrfs snapshot (same as btrbk)
- **Backup**: A logical group of snapshots taken at the same point in time across multiple volumes. For example, all snapshots created at `20250909T2343` for volumes `@`, `@home`, and `@srv` form one "backup". This is NOT the same as btrbk's remote backup feature.

<!-- Note: Not all snapshots in a backup group are necessarily created at the exact same time. If a snapshot is missing for a volume at the backup timestamp, the group will include the most recent previous snapshot of that volume instead. -->

## Configuration

The tool can optionally use a [config.yaml](config.yaml) file to configure its behavior. By default, it looks for this file in the current working directory.

**Note:** The config file is entirely optional. If not provided, the tool will automatically detect configuration by parsing the output of `btrbk list`.

### Configuration File Options

```yaml
btrbk_config_file: ./btrbk.conf
old_format: "{{.SubvolName}}.old"
default_subvols_restore_list:
    - "@"

# Path configurations
paths:
    snaps: ./mnt/@snaps
    target: ./mnt
    meta: ./mnt/@snaps/.meta
    meta_trash: ./mnt/@snaps/.meta/.trash
```

**All fields in the config file are optional.** The tool will use sensible defaults or auto-detect values based on your btrbk configuration.

### Configuration Details

- **btrbk_config_file**: Path to your btrbk configuration file. This is used to enumerate available snapshots via `btrbk list`.

- **default_subvols_restore_list**: Specifies which subvolumes should be selected by default in the restore interface. Useful for quickly restoring common volumes like the root filesystem (`@`).

- **old_format**: Template string that determines how existing subvolumes are renamed during restore operations. Uses Go template syntax. Available variables: `SubvolName`, `Timestamp`.

- **paths.snaps**: The directory containing your btrbk snapshots.

- **paths.target**: The directory where subvolumes will be restored to.

- **paths.meta**: Override the default metadata directory location. Defaults to `.meta` inside the snapshots directory.

- **paths.meta_trash**: Override the default trash directory location. Defaults to `.trash` inside the metadata directory.

## Features

### Deleting Snaps

### Protecting Snaps

The tool creates the following metadata structure in the snapshots directory:
```.meta
└── 20250909T2343
    ├── info.yaml
    ├── @
    ├── @home
    └── @srv
```

When unchecking a protected backup, it is moved to the `.trash` directory instead of being immediately deleted.

`info.yaml` contains:
- User notes for the backup
- Restoration dates
- Tags

<!-- ### Unified Snaps

Manager may group snaps by date.

if `snapshot_create onchange` option is applied, -->
