package snaps

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func (mng *BackupManager) Protect(timestamp string, note ProtectionNote) error {
	var err error
	if mng.isInTrash(timestamp) {
		err = mng.restoreFromTrash(timestamp)
	} else if !mng.isProtected(timestamp) {
		err = mng.persistBackup(timestamp)
	}
	if err != nil {
		return err
	}

	return mng.attachNote(timestamp, note)
}

func (mng *BackupManager) ClearCache() {
	mng.collectResult = nil
}

func (mng *BackupManager) FreePersistance(timestamp string) error {
	metaTimestampDir := filepath.Join(mng.metaDir, timestamp)
	trashDir := filepath.Join(mng.trashDir, timestamp)

	if err := os.MkdirAll(filepath.Dir(trashDir), 0755); err != nil {
		return fmt.Errorf("failed to create .trash directory: %w", err)
	}

	if err := os.Rename(metaTimestampDir, trashDir); err != nil {
		return fmt.Errorf("failed to move %s to trash: %w", timestamp, err)
	}

	return nil
}

func (mng *BackupManager) isProtected(timestamp string) bool {
	metaTimestampDir := filepath.Join(mng.metaDir, timestamp)
	_, err := os.Stat(metaTimestampDir)
	return err == nil
}

func (mng *BackupManager) isInTrash(timestamp string) bool {
	trashDir := filepath.Join(mng.trashDir, timestamp)
	_, err := os.Stat(trashDir)
	return err == nil
}

func (mng *BackupManager) restoreFromTrash(timestamp string) error {
	trashDir := filepath.Join(mng.trashDir, timestamp)
	metaTimestampDir := filepath.Join(mng.metaDir, timestamp)

	if err := os.Rename(trashDir, metaTimestampDir); err != nil {
		return fmt.Errorf("failed to restore %s from trash: %w", timestamp, err)
	}

	return nil
}

func (mng *BackupManager) persistBackup(timestamp string) error {
	errSetupMeta := mng.setupSubvolumeMeta(timestamp)
	snapCollection, errCollect := mng.Collect()

	err := errors.Join(errSetupMeta, errCollect)
	if err != nil {
		return err
	}

	metaTimestampDir := filepath.Join(mng.metaDir, timestamp)

	if len(snapCollection.SubvolNames) == 0 {
		return fmt.Errorf("no backups found to protect at timestamp %s", timestamp)
	}

	for _, subvolName := range snapCollection.SubvolNames {
		srcPath := filepath.Join(mng.dir, subvolName+"."+timestamp)
		dstPath := filepath.Join(metaTimestampDir, subvolName)

		var stderr bytes.Buffer
		cmd := exec.Command("btrfs", "subvolume", "snapshot", "-r", srcPath, dstPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to snapshot %s: %w", subvolName, err)
		}
		if stderr.Len() > 0 {
			return errors.New(stderr.String())
		}
	}

	return nil
}

func (mng *BackupManager) attachNote(timestamp string, note ProtectionNote) error {
	yamlData, err := yaml.Marshal(&note)
	if err != nil {
		return fmt.Errorf("failed to marshal note to YAML: %w", err)
	}

	infoPath := filepath.Join(mng.metaDir, timestamp, "info.yaml")
	if err := os.WriteFile(infoPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write info.yaml: %w", err)
	}

	return nil
}
