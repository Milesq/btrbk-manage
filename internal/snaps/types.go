package snaps

import "fmt"

type Snapshot struct {
	Timestamp, SubvolName string
}

type Backup struct {
	Timestamp      string
	Items          []Snapshot
	IsProtected    bool
	ProtectionNote ProtectionNote
}

func (b Backup) String() string {
	return fmt.Sprintf("Backup{Timestamp: %s, Items: %d, IsProtected: %t}",
		b.Timestamp, len(b.Items), b.IsProtected)
}

type CollectResult struct {
	Backups                  []Backup
	SubvolNames              []string
	ProtectedN, UnprotectedN int
}
