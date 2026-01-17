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

	timestampsToSkipMap := make(map[string]struct{})
	backups, err := collectBackups(filepath.Join(mng.paths.Meta), &timestampsToSkipMap, Backup{IsProtected: true})
	trashed, errTrash := collectBackups(filepath.Join(mng.paths.MetaTrash), &timestampsToSkipMap, Backup{IsTrashed: true})
	unprotected, errSnaps := collectSnapshots(mng.paths.Snaps, timestampsToSkipMap)

	if err = errors.Join(err, errTrash, errSnaps); err != nil {
		return CollectResult{}, err
	}

	allBackups := append(backups, trashed...)
	allBackups = append(allBackups, unprotected.Backups...)

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
