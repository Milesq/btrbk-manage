package snaps

type Snapshot struct {
	Path      string
	BaseName  string
	Timestamp string
}

type Group struct {
	Timestamp string
	Items     []Snapshot
}
