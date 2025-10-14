package btrbk

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type ListResult struct {
	Source       string
	SnapPath     string
	SnapshotName string
}

func List(configFile string) ([]ListResult, error) {
	program := "btrbk"
	args := []string{"list", "source"}

	if configFile != "" {
		args = append(args, "-c", configFile)
	}

	cmd := exec.Command(program, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("btrbk list failed: %w\nstderr: %s", err, stderr.String())
	}

	return parseListOutput(stdout.String())
}

func parseListOutput(output string) ([]ListResult, error) {
	var results []ListResult

	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		result := ListResult{
			Source:       fields[0],
			SnapPath:     fields[1],
			SnapshotName: fields[2],
		}

		results = append(results, result)
	}

	return results, nil
}
