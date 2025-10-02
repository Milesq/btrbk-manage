package snaps

type Snapshot struct {
	Timestamp, SubvolName, BaseName string
}

type Group struct {
	Timestamp string
	Items     []Snapshot
}
