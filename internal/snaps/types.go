package snaps

type Snapshot struct {
	Timestamp, SubvolName string
}

type Backup struct {
	Timestamp      string
	Items          []Snapshot
	IsProtected    bool
	ProtectionNote ProtectionNote
}

type CollectResult struct {
	Backups                  []Backup
	SubvolNames              []string
	ProtectedN, UnprotectedN int
}
