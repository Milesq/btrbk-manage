package snaps

import "fmt"

func (mng *BackupManager) Restore(backup Backup, subvolumes []string) error {
	fmt.Println("restore", backup.Items)
	return nil
}
