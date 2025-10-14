package snaps

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func GetManagerForDirectory(snapDirectory, metaDir, trashDir string) BackupManager {
	return BackupManager{snapDirectory, metaDir, trashDir, nil, nil, nil}
}

type BackupManager struct {
	dir      string
	metaDir  string
	trashDir string

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
	metaDir := filepath.Join(mng.metaDir, timestamp)
	return os.MkdirAll(metaDir, 0755)
}
