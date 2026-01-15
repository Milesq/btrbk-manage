package snaps

import (
	"fmt"
	"os"
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
