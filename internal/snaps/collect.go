package snaps

import (
	"errors"
	"path/filepath"
	"sort"
)

func (mng *BackupManager) Collect() (CollectResult, error) {
	if mng.collectResult != nil {
		return *mng.collectResult, nil
	}

	backups, timestampBackupMap, err := collectBackups(filepath.Join(mng.dir, ".meta"))
	unprotected, errSnaps := collectSnapshots(mng.dir, timestampBackupMap)

	if err = errors.Join(err, errSnaps); err != nil {
		return CollectResult{}, err
	}

	allBackups := append(backups, unprotected.Backups...)

	sort.Slice(allBackups, func(i, j int) bool {
		return allBackups[i].Timestamp > allBackups[j].Timestamp
	})

	result := CollectResult{
		Backups:      allBackups,
		SubvolNames:  unprotected.SubvolNames,
		ProtectedN:   len(backups),
		UnprotectedN: len(unprotected.Backups),
	}

	mng.collectResult = &result

	return result, err
}
