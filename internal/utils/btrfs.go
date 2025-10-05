package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"milesq.dev/btrbk-manage/internal"
)

func BtrfsDelete(subvolPath string) error {
	var stderr bytes.Buffer
	program := "btrfs"
	args := []string{"subvolume", "delete", subvolPath}

	if internal.Env == "dev" {
		program = "sudo"
		args = append([]string{"btrfs"}, args...)
	}

	cmd := exec.Command(program, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete btrfs subvolume %s: %w, stderr: %s", subvolPath, err, stderr.String())
	}

	return nil
}
