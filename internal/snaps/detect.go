package snaps

import (
	"regexp"
)

var timestampRe = regexp.MustCompile(`(@.*)\.(\d{8}T\d{4})$`)

func detectSnapshot(base string) (string, string, bool) {
	m := timestampRe.FindStringSubmatch(base)

	if m == nil {
		return "", "", false
	}
	return m[1], m[2], true
}
