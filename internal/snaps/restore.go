package snaps

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

func (mng *BackupManager) Restore(backup Backup, subvolumes []string) error {
	if backup.IsProtected {
		if err := mng.addRestorationDate(backup.Timestamp); err != nil {
			return fmt.Errorf("failed to update restoration date: %w", err)
		}
	}

	if err := mng.executeHook("post-restore"); err != nil {
		return fmt.Errorf("post-restore hook failed: %w", err)
	}

	return nil
}

func (mng *BackupManager) executeHook(name string) error {
	hookPath := filepath.Join(mng.hooksDir, name+".sh")

	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return nil
	}

	var stderr bytes.Buffer
	cmd := exec.Command(hookPath)
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(), fmt.Sprintf("RESTORE_PATH=%s", mng.dir))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to execute %s: %w, stderr: %s", hookPath, err, stderr.String())
	}

	return nil
}

func (mng *BackupManager) addRestorationDate(timestamp string) error {
	infoPath := filepath.Join(mng.metaDir, timestamp, "info.yaml")

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
