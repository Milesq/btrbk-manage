package snaps

func GetManagerForDirectory(snapDirectory string) BackupManager {
	return BackupManager{snapDirectory}
}

type BackupManager struct {
	dir string
}
