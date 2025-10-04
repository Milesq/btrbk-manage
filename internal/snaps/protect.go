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
	if mng.isProtected(timestamp) {
		return fmt.Errorf("backup at timestamp %s is already protected", timestamp)
	}

	err := mng.persistBackup(timestamp)
	if err != nil {
		return err
	}

	return mng.attachNote(timestamp, note)
}

func (mng *BackupManager) isProtected(timestamp string) bool {
	metaTimestampDir := filepath.Join(mng.dir, ".meta", timestamp)
	_, err := os.Stat(metaTimestampDir)
	return err == nil
}

func (mng *BackupManager) persistBackup(timestamp string) error {
	errSetupMeta := mng.setupSubvolumeMeta(timestamp)
	backups, errCollect := mng.Collect()

	err := errors.Join(errSetupMeta, errCollect)
	if err != nil {
		return err
	}

	metaTimestampDir := filepath.Join(mng.dir, ".meta", timestamp)

	if len(backups.SubvolNames) == 0 {
		return fmt.Errorf("no backups found to protect at timestamp %s", timestamp)
	}

	for _, subvolName := range backups.SubvolNames {
		srcPath := filepath.Join(mng.dir, subvolName+"."+timestamp)
		dstPath := filepath.Join(metaTimestampDir, subvolName)

		var stderr bytes.Buffer
		cmd := exec.Command("btrfs", "subvolume", "snapshot", "-r", srcPath, dstPath)
		cmd.Stdout = os.Stdout
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

	infoPath := filepath.Join(mng.dir, ".meta", timestamp, "info.yaml")
	if err := os.WriteFile(infoPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write info.yaml: %w", err)
	}

	return nil
}
