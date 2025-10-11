package snaps

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func collectBackups(dir string) ([]Backup, map[string]struct{}, error) {
	subvolNamesMap := make(map[string]struct{})

	metaDir, err := os.ReadDir(dir)
	if err != nil {
		return []Backup{}, nil, err
	}

	backups := make([]Backup, 0, len(metaDir))

	var errs []error
	for _, entry := range metaDir {
		if !entry.IsDir() || entry.Name() == ".trash" {
			continue
		}

		timestamp := entry.Name()
		backupPath := filepath.Join(dir, timestamp)

		snapLinks, err := os.ReadDir(backupPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		subvolNamesMap[timestamp] = struct{}{}

		items := make([]Snapshot, 0)
		for _, link := range snapLinks {
			if link.Name() == "info.yaml" {
				continue
			}

			items = append(items, Snapshot{
				SubvolName: link.Name(),
				Timestamp:  timestamp,
			})
		}

		if len(items) == 0 {
			continue
		}

		var protectionNote ProtectionNote
		infoPath := filepath.Join(backupPath, "info.yaml")
		if data, err := os.ReadFile(infoPath); err == nil {
			_ = yaml.Unmarshal(data, &protectionNote)
		}

		backups = append(backups, Backup{
			Timestamp:      timestamp,
			Items:          items,
			IsProtected:    true,
			ProtectionNote: protectionNote,
		})
	}

	return backups, subvolNamesMap, errors.Join(errs...)
}
