package snaps

import (
	"os"
	"sort"
)

func collectSnapshots(dir string) (CollectResult, error) {
	gmap := make(map[string][]Snapshot)
	subvolNamesMap := make(map[string]struct{})

	snapDir, err := os.ReadDir(dir)
	if err != nil {
		return CollectResult{}, err
	}

	for _, v := range snapDir {
		if !v.IsDir() {
			continue
		}
		name := v.Name()
		subvolName, snapTimeStmp, ok := detectSnapshot(name)

		if !ok {
			continue
		}

		subvolNamesMap[subvolName] = struct{}{}
		gmap[snapTimeStmp] = append(gmap[snapTimeStmp], Snapshot{
			SubvolName: subvolName,
			Timestamp:  snapTimeStmp,
		})
	}

	if len(gmap) == 0 {
		return CollectResult{}, nil
	}

	backups := make([]Backup, 0, len(gmap))
	for ts, items := range gmap {
		sort.Slice(items, func(i, j int) bool { return items[i].SubvolName < items[j].SubvolName })
		backups = append(backups, Backup{Timestamp: ts, Items: items})
	}

	sort.Slice(backups, func(i, j int) bool { return backups[i].Timestamp > backups[j].Timestamp })

	subvolNames := make([]string, 0, len(subvolNamesMap))
	for name := range subvolNamesMap {
		subvolNames = append(subvolNames, name)
	}
	sort.Strings(subvolNames)

	return CollectResult{
		Backups:     backups,
		SubvolNames: subvolNames,
	}, nil

}
