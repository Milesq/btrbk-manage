package snaps

type Snapshot struct {
	Timestamp, SubvolName, BaseName string
}

type Backup struct {
	Timestamp      string
	Items          []Snapshot
	IsProtected    bool
	ProtectionNote ProtectionNote
}

type CollectResult struct {
	Backups     []Backup
	SubvolNames []string
	TotalCount  int
}
