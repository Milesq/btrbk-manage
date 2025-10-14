package snaps

import "fmt"

func (mng *BackupManager) Restore(backup Backup) error {
	fmt.Println("restore", backup.Items)
	return nil
}
