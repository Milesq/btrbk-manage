package snaps

import (
	"fmt"
	"strings"
)

type ProtectionNote struct {
	Note             string   `yaml:"note,omitempty"`
	Reason           string   `yaml:"reason,omitempty"`
	Tags             []string `yaml:"tags,omitempty"`
	RestorationDates []string `yaml:"restoration_dates,omitempty"`
}

const building_params_n int = 3

func GetProtectionNote(values []string) (note ProtectionNote, err error) {
	if len(values) != building_params_n {
		return ProtectionNote{}, fmt.Errorf("building protection note failed. Expected %d params, got %d", building_params_n, len(values))
	}
	note.Note = values[0]
	note.Reason = values[1]
	note.Tags = strings.Split(values[2], ",")
	for j := range note.Tags {
		note.Tags[j] = strings.TrimSpace(note.Tags[j])
	}
	return
}
