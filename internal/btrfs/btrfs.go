package btrfs

import (
	"bytes"
	"fmt"
	"os/exec"

	"milesq.dev/btrbk-manage/internal"
)

func SubvolDelete(subvolPath string) error {
	var stderr bytes.Buffer
	program := "btrfs"
	args := []string{"subvolume", "delete", subvolPath}

	if internal.Env == "dev" {
		program = "sudo"
		args = append([]string{"btrfs"}, args...)
	}

	cmd := exec.Command(program, args...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete btrfs subvolume %s: %w, stderr: %s", subvolPath, err, stderr.String())
	}

	return nil
}

func Snapshot(source, dest string, ro bool) error {
	var stderr bytes.Buffer
	program := "btrfs"
	args := []string{"subvolume", "snapshot"}

	if ro {
		args = append(args, "-r")
	}

	args = append(args, source, dest)

	cmd := exec.Command(program, args...)
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create btrfs snapshot from %s to %s: %w, stderr: %s", source, dest, err, stderr.String())
	}

	return nil
}
