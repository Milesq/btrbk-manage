package snaps

import (
	"fmt"
	"os"
	"path/filepath"

	"milesq.dev/btrbk-manage/internal/utils"
)

func (mng *BackupManager) RemoveFromTrash(timestamp string) error {
	trashDir := filepath.Join(mng.dir, ".meta/.trash", timestamp)

	entries, err := os.ReadDir(trashDir)
	if err != nil {
		return fmt.Errorf("failed to read trash directory %s: %w", timestamp, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subvolPath := filepath.Join(trashDir, entry.Name())
			if err := utils.BtrfsDelete(subvolPath); err != nil {
				return err
			}
		}
	}

	if err := os.RemoveAll(trashDir); err != nil {
		return fmt.Errorf("failed to remove trash directory %s: %w", timestamp, err)
	}

	return nil
}
