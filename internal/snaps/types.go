package snaps

type Snapshot struct {
	Timestamp, BaseName string
}

type Group struct {
	Timestamp string
	Items     []Snapshot
}
