package snaps

import (
	"fmt"
	"os"
	"path/filepath"

	"milesq.dev/btrbk-manage/internal/utils"
)

func (mng *BackupManager) Delete(backup Backup) error {
	if err := mng.DeleteVirtualBackup(backup); err != nil {
		return err
	}

	if !backup.IsProtected && !backup.IsTrashed {
		return nil
	}

	dir := mng.metaDir
	if backup.IsTrashed {
		dir = mng.trashDir
	}

	backupDir := filepath.Join(dir, backup.Timestamp)

	return mng.DeletePhysicalBackup(backupDir)
}

func (mng *BackupManager) DeleteVirtualBackup(backup Backup) error {
	for _, item := range backup.Items {
		subvolPath := filepath.Join(mng.dir, item.SubvolName+"."+item.Timestamp)
		if err := utils.BtrfsDelete(subvolPath); err != nil {
			return fmt.Errorf("failed to delete virtual backup %s: %w", item.SubvolName, err)
		}
	}
	return nil
}

func (mng *BackupManager) DeletePhysicalBackup(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subvolPath := filepath.Join(dir, entry.Name())
			if err := utils.BtrfsDelete(subvolPath); err != nil {
				return fmt.Errorf("failed to delete subvolume %s: %w", subvolPath, err)
			}
		}
	}

	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("failed to remove directory %s: %w", dir, err)
	}

	return nil
}

func (mng *BackupManager) RemoveFromTrash(backup Backup) error {
	if !backup.IsTrashed {
		return fmt.Errorf("backup %s is not in trash", backup.Timestamp)
	}

	trashDir := filepath.Join(mng.trashDir, backup.Timestamp)
	return mng.DeletePhysicalBackup(trashDir)
}
