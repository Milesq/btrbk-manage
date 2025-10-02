package utils

import (
	"testing"
)

func TestMinMax(t *testing.T) {
	tests := []struct {
		minInt int
		n      int
		maxInt int
		want   int
	}{
		{1, 2, 3, 2},
		{1, -1, 3, 1},
		{-5, 0, 5, 0},
		{-10, -5, 5, -5},
	}

	for _, tt := range tests {
		got := MinMax(tt.minInt, tt.n, tt.maxInt)
		if got != tt.want {
			t.Errorf("MinMax(%d, %d, %d) = %d, want %d", tt.minInt, tt.n, tt.maxInt, got, tt.want)
		}
	}
}
