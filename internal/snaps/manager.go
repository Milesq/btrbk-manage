package snaps

import (
	"fmt"
	"slices"
)

func GetManagerForDirectory(snapDirectory string) BackupManager {
	return BackupManager{snapDirectory, nil, nil}
}

type BackupManager struct {
	dir                 string
	availableSubvolumes []string
	subvolumes          []string
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
