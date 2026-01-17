package snaps

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"milesq.dev/btrbk-manage/internal/app"
)

func GetBackupManager(paths app.Paths) BackupManager {
	return BackupManager{paths: paths}
}

type BackupManager struct {
	paths app.Paths

	availableSubvolumes []string
	subvolumes          []string
	collectResult       *CollectResult
}

func (mng *BackupManager) AvailableSubvolumes() []string {
	result := make([]string, len(mng.availableSubvolumes))
	copy(result, mng.availableSubvolumes)
	return result
}

func (mng *BackupManager) SetSubvolumes(subvolumes []string) error {
	for _, subvol := range subvolumes {
		if !mng.isAvailable(subvol) {
			return fmt.Errorf("subvolume %q is not available", subvol)
		}
	}
	mng.subvolumes = subvolumes
	return nil
}

func (mng *BackupManager) isAvailable(subvolume string) bool {
	return slices.Contains(mng.availableSubvolumes, subvolume)
}

func (mng *BackupManager) setupSubvolumeMeta(timestamp string) error {
	metaDir := filepath.Join(mng.paths.Meta, timestamp)
	return os.MkdirAll(metaDir, 0755)
}
