package snaps

func (mng *BackupManager) Collect() (CollectResult, error) {
	if mng.collectResult != nil {
		return *mng.collectResult, nil
	}

	unprotected, err := collectSnapshots(mng.dir)
	mng.collectResult = &unprotected

	collectBackups(mng.dir)
	return unprotected, err
}
