package snaps

import "testing"

func TestDetectSnapshot(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		ok       bool
	}{
		{
			name:     "valid snapshot with standard timestamp",
			input:    "backup@home.20240928T1430",
			expected: "20240928T1430",
			ok:       true,
		},
		{
			name:     "valid snapshot with different name",
			input:    "mydata@root.20231215T0900",
			expected: "20231215T0900",
			ok:       true,
		},
		{
			name:     "invalid format - no timestamp",
			input:    "backup@home",
			expected: "",
			ok:       false,
		},
		{
			name:     "invalid format - wrong timestamp format",
			input:    "backup@home.20240928",
			expected: "",
			ok:       false,
		},
		{
			name:     "invalid format - no @ symbol",
			input:    "backup.20240928T1430",
			expected: "",
			ok:       false,
		},
		{
			name:     "invalid format - extra characters after timestamp",
			input:    "backup@home.20240928T1430extra",
			expected: "",
			ok:       false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			ok:       false,
		},
		{
			name:     "just timestamp without prefix",
			input:    "20240928T1430",
			expected: "",
			ok:       false,
		},
		{
			name:     "valid with longer prefix",
			input:    "very_long_backup_name@some_subvol.20240928T1430",
			expected: "20240928T1430",
			ok:       true,
		},
		{
			name:     "valid with numbers in prefix",
			input:    "backup123@home456.20240928T1430",
			expected: "20240928T1430",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := detectSnapshot(tt.input)
			if ok != tt.ok {
				t.Errorf("detectSnapshot(%q) ok = %v, want %v", tt.input, ok, tt.ok)
			}
			if result != tt.expected {
				t.Errorf("detectSnapshot(%q) result = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
