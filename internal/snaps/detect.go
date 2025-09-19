package snaps

import (
	"regexp"
)

var tsRe = regexp.MustCompile(`@.*\.(\d{8}T\d{4})$`)

func detectSnapshot(base string) (string, bool) {
	m := tsRe.FindStringSubmatch(base)

	if m == nil {
		return "", false
	}
	return m[1], true
}
