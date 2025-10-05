package snaps

import (
	"os"
	"sort"
)

func (mng *BackupManager) Collect() (CollectResult, error) {
	if mng.collectResult != nil {
		return *mng.collectResult, nil
	}
	gmap := make(map[string][]Snapshot)
	subvolNamesMap := make(map[string]struct{})

	snapDir, err := os.ReadDir(mng.dir)
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
			BaseName:   name,
			SubvolName: subvolName,
			Timestamp:  snapTimeStmp,
		})
	}

	if len(gmap) == 0 {
		return CollectResult{}, nil
	}

	groups := make([]Group, 0, len(gmap))
	for ts, items := range gmap {
		sort.Slice(items, func(i, j int) bool { return items[i].BaseName < items[j].BaseName })
		groups = append(groups, Group{Timestamp: ts, Items: items, IsProtected: mng.isProtected(ts)})
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i].Timestamp > groups[j].Timestamp })

	subvolNames := make([]string, 0, len(subvolNamesMap))
	for name := range subvolNamesMap {
		subvolNames = append(subvolNames, name)
	}
	sort.Strings(subvolNames)

	mng.collectResult = &CollectResult{
		Groups:      groups,
		SubvolNames: subvolNames,
		TotalCount:  len(snapDir),
	}

	return *mng.collectResult, nil
}
