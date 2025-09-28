package snaps

import (
	"os"
	"sort"
)

func Collect(dir string) ([]Group, int, error) {
	gmap := make(map[string][]Snapshot)

	snapDir, err := os.ReadDir(dir)
	if err != nil {
		return nil, 0, nil
	}

	for _, v := range snapDir {
		if !v.IsDir() {
			continue
		}
		name := v.Name()
		snapTimeStmp, ok := detectSnapshot(name)

		if !ok {
			continue
		}

		gmap[snapTimeStmp] = append(gmap[snapTimeStmp], Snapshot{
			BaseName:  name,
			Timestamp: snapTimeStmp,
		})
	}

	if len(gmap) == 0 {
		return nil, 0, err
	}

	groups := make([]Group, 0, len(gmap))
	for ts, items := range gmap {
		sort.Slice(items, func(i, j int) bool { return items[i].BaseName < items[j].BaseName })
		groups = append(groups, Group{Timestamp: ts, Items: items})
	}

	sort.Slice(groups, func(i, j int) bool { return groups[i].Timestamp > groups[j].Timestamp })

	return groups, len(snapDir), nil
}
