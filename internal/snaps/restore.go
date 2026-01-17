package snaps

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
	"milesq.dev/btrbk-manage/internal/btrfs"
)

func (mng *BackupManager) Restore(backup Backup, subvolumes []string) error {
	subvolSet := make(map[string]struct{})
	for _, sv := range subvolumes {
		subvolSet[sv] = struct{}{}
	}

	for _, snapshot := range backup.Items {
		if _, ok := subvolSet[snapshot.SubvolName]; len(subvolumes) > 0 && !ok {
			continue
		}

		if err := mng.restoreSnapshot(snapshot); err != nil {
			return fmt.Errorf("failed to restore %s: %w", snapshot.SubvolName, err)
		}
	}

	if backup.IsProtected {
		if err := mng.addRestorationDate(backup.Timestamp); err != nil {
			return fmt.Errorf("failed to update restoration date: %w", err)
		}
	}

	env := append(os.Environ(), fmt.Sprintf("SNAPSHOT=%s", filepath.Join(mng.paths.Snaps)))

	if err := mng.executeHook("post-restore", env); err != nil {
		return fmt.Errorf("post-restore hook failed: %w", err)
	}

	return nil
}

func (mng *BackupManager) restoreSnapshot(snapshot Snapshot) error {
	sourcePath, err := mng.getSnapshotPath(snapshot)

	if err != nil {
		return fmt.Errorf("failed to get snapshot path: %w", err)
	}

	targetPath := filepath.Join(mng.paths.Target, snapshot.SubvolName)

	if _, err := os.Stat(targetPath); err == nil {
		oldPath := targetPath + ".old"

		if _, err := os.Stat(oldPath); err == nil {
			if err := btrfs.SubvolDelete(oldPath); err != nil {
				return fmt.Errorf("failed to delete existing .old subvolume: %w", err)
			}
		}

		log.Printf("Renaming existing subvolume %s to %s", targetPath, oldPath)
		if err := os.Rename(targetPath, oldPath); err != nil {
			return fmt.Errorf("failed to rename %s to %s: %w", targetPath, oldPath, err)
		}
	}

	if err := btrfs.Snapshot(sourcePath, targetPath, false); err != nil {
		return err
	}

	return nil
}

func (mng *BackupManager) getSnapshotPath(snapshot Snapshot) (string, error) {
	sourcePath := filepath.Join(mng.paths.Snaps, snapshot.SubvolName+"."+snapshot.Timestamp)
	var err, err2 error

	if _, err = os.Stat(sourcePath); err == nil {
		return sourcePath, nil
	}

	sourcePath = filepath.Join(mng.paths.Meta, snapshot.Timestamp, snapshot.SubvolName)
	if _, err2 = os.Stat(sourcePath); err2 == nil {
		return sourcePath, nil
	}

	return "", errors.Join(errors.New("Cannot determine snapshot path"), err, err2)
}

func (mng *BackupManager) executeHook(name string, env []string) error {
	hookPath := filepath.Join(mng.paths.Hooks, name+".sh")

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return nil
	}

	var stderr bytes.Buffer
	cmd := exec.Command(hookPath)
	cmd.Stderr = &stderr
	cmd.Env = env

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute %s: %w, stderr: %s", hookPath, err, stderr.String())
	}

	return nil
}

func (mng *BackupManager) addRestorationDate(timestamp string) error {
	infoPath := filepath.Join(mng.paths.Meta, timestamp, "info.yaml")

	var note ProtectionNote
	if data, err := os.ReadFile(infoPath); err == nil {
		_ = yaml.Unmarshal(data, &note)
	}

	currentDate := time.Now().String()
	note.RestorationDates = append(note.RestorationDates, currentDate)

	yamlData, err := yaml.Marshal(&note)
	if err != nil {
		return fmt.Errorf("failed to marshal note to YAML: %w", err)
	}

	if err := os.WriteFile(infoPath, yamlData, 0644); err != nil {
		return fmt.Errorf("failed to write info.yaml: %w", err)
	}

	return nil
}
