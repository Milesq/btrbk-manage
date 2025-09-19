package snaps

import (
	"errors"
	"io/fs"
	"path/filepath"
	"sort"
)

func Collect(dir string) ([]Group, int, error) {
	gmap := make(map[string][]Snapshot)
	var total int

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Skip unreadable entries but continue.
			return nil
		}
		// Only consider directories (typical for btrfs snapshots), but allow files too in case of send streams, etc.
		base := filepath.Base(path)
		ts, ok := detectSnapshot(base)
		if !ok {
			return nil
		}
		// Optional: avoid descending into nested snapshots too deep—WalkDir already handles recursion.
		info, e := d.Info()
		if e == nil && info.Mode().IsDir() || (e == nil && info.Mode().IsRegular()) {
			gmap[ts] = append(gmap[ts], Snapshot{
				Path:      path,
				BaseName:  base,
				Timestamp: ts,
			})
			total++
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}

	if len(gmap) == 0 {
		return nil, 0, errors.New("no snapshots matching pattern *.YYYYMMDDTHHMMSS found")
	}

	// Build groups slice and sort members by name.
	groups := make([]Group, 0, len(gmap))
	for ts, items := range gmap {
		sort.Slice(items, func(i, j int) bool { return items[i].BaseName < items[j].BaseName })
		groups = append(groups, Group{Timestamp: ts, Items: items})
	}
	// Sort groups by timestamp desc (string compare works for this format).
	sort.Slice(groups, func(i, j int) bool { return groups[i].Timestamp > groups[j].Timestamp })

	return groups, total, nil
}
