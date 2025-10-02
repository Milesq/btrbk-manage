package snaps

import "testing"

func TestDetectSnapshot(t *testing.T) {
	tests := []struct {
		name              string
		input             string
		expectedSubvol    string
		expectedTimestamp string
		ok                bool
	}{
		{
			name:              "valid snapshot with standard timestamp",
			input:             "backup@home.20240928T1430",
			expectedSubvol:    "@home",
			expectedTimestamp: "20240928T1430",
			ok:                true,
		},
		{
			name:              "valid snapshot with different name",
			input:             "mydata@root.20231215T0900",
			expectedSubvol:    "@root",
			expectedTimestamp: "20231215T0900",
			ok:                true,
		},
		{
			name:              "invalid format - no timestamp",
			input:             "backup@home",
			expectedSubvol:    "",
			expectedTimestamp: "",
			ok:                false,
		},
		{
			name:              "invalid format - wrong timestamp format",
			input:             "backup@home.20240928",
			expectedSubvol:    "",
			expectedTimestamp: "",
			ok:                false,
		},
		{
			name:              "invalid format - no @ symbol",
			input:             "backup.20240928T1430",
			expectedSubvol:    "",
			expectedTimestamp: "",
			ok:                false,
		},
		{
			name:              "invalid format - extra characters after timestamp",
			input:             "backup@home.20240928T1430extra",
			expectedSubvol:    "",
			expectedTimestamp: "",
			ok:                false,
		},
		{
			name:              "empty string",
			input:             "",
			expectedSubvol:    "",
			expectedTimestamp: "",
			ok:                false,
		},
		{
			name:              "just timestamp without prefix",
			input:             "20240928T1430",
			expectedSubvol:    "",
			expectedTimestamp: "",
			ok:                false,
		},
		{
			name:              "valid with longer prefix",
			input:             "very_long_backup_name@some_subvol.20240928T1430",
			expectedSubvol:    "@some_subvol",
			expectedTimestamp: "20240928T1430",
			ok:                true,
		},
		{
			name:              "valid with numbers in prefix",
			input:             "backup123@home456.20240928T1430",
			expectedSubvol:    "@home456",
			expectedTimestamp: "20240928T1430",
			ok:                true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subvol, timestamp, ok := detectSnapshot(tt.input)
			if ok != tt.ok {
				t.Errorf("detectSnapshot(%q) ok = %v, want %v", tt.input, ok, tt.ok)
			}
			if subvol != tt.expectedSubvol {
				t.Errorf("detectSnapshot(%q) subvol = %q, want %q", tt.input, subvol, tt.expectedSubvol)
			}
			if timestamp != tt.expectedTimestamp {
				t.Errorf("detectSnapshot(%q) timestamp = %q, want %q", tt.input, timestamp, tt.expectedTimestamp)
			}
		})
	}
}
