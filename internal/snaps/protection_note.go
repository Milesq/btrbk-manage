package snaps

type ProtectionNote struct {
	Note             string   `yaml:"note,omitempty"`
	Reason           string   `yaml:"reason,omitempty"`
	Tags             []string `yaml:"tags,omitempty"`
	RestorationDates []string `yaml:"restoration_dates,omitempty"`
}
