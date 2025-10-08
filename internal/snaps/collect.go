package snaps

func (mng *BackupManager) Collect() (CollectResult, error) {
	if mng.collectResult != nil {
		return *mng.collectResult, nil
	}

	unprotected, err := mng.collectSnapshotsFrom(mng.dir)
	mng.collectResult = &unprotected
	return unprotected, err
}
