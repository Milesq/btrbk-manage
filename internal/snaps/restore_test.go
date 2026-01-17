package snaps

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
	"milesq.dev/btrbk-manage/internal/app"
)

func setupRestoreTestDirs(t *testing.T) (tempDir string, paths app.Paths, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "restore_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	paths = app.Paths{
		Snaps:     filepath.Join(tempDir, "snapshots"),
		Target:    filepath.Join(tempDir, "target"),
		Meta:      filepath.Join(tempDir, "meta"),
		MetaTrash: filepath.Join(tempDir, "trash"),
		Hooks:     filepath.Join(tempDir, "hooks"),
	}

	for _, dir := range []string{paths.Snaps, paths.Meta, paths.MetaTrash, paths.Hooks, paths.Target} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			os.RemoveAll(tempDir)
			t.Fatalf("Failed to create dir %s: %v", dir, err)
		}
	}

	cleanup = func() {
		os.RemoveAll(tempDir)
	}

	return
}

func createHookScript(t *testing.T, hooksDir, name, content string) {
	t.Helper()
	hookPath := filepath.Join(hooksDir, name+".sh")
	if err := os.WriteFile(hookPath, []byte("#!/bin/bash\n"+content), 0755); err != nil {
		t.Fatalf("Failed to create hook script %s: %v", name, err)
	}
}

func TestRestore_FiltersBySubvolumes(t *testing.T) {
	_, paths, cleanup := setupRestoreTestDirs(t)
	defer cleanup()

	mng := GetBackupManager(paths)

	timestamp := "20240928T1430"
	timestampMetaDir := filepath.Join(paths.Meta, timestamp)
	if err := os.MkdirAll(timestampMetaDir, 0755); err != nil {
		t.Fatalf("Failed to create timestamp meta dir: %v", err)
	}

	// Create source snapshot directories (only for filtered subvolumes)
	// @root is NOT created to verify it's filtered out
	for _, subvol := range []string{"@home", "@var"} {
		subvolDir := filepath.Join(timestampMetaDir, subvol)
		if err := os.MkdirAll(subvolDir, 0755); err != nil {
			t.Fatalf("Failed to create subvol dir %s: %v", subvol, err)
		}
	}

	backup := Backup{
		Timestamp:   timestamp,
		IsProtected: false,
		Items: []Snapshot{
			{Timestamp: timestamp, SubvolName: "@root"},
			{Timestamp: timestamp, SubvolName: "@home"},
			{Timestamp: timestamp, SubvolName: "@var"},
		},
	}

	// This will fail at btrfs snapshot command but verifies filtering works
	// (if @root wasn't filtered, it would fail with "source snapshot does not exist")
	err := mng.Restore(backup, []string{"@home", "@var"})
	// Expect error from btrfs command (not available in test env), but NOT
	// "source snapshot does not exist" for @root (which would mean filtering failed)
	if err != nil && strings.Contains(err.Error(), "failed to restore @root") {
		t.Errorf("Restore() should have filtered out @root, but it wasn't filtered")
	}
}

func TestRestore_AddsRestorationDateForProtectedBackup(t *testing.T) {
	_, paths, cleanup := setupRestoreTestDirs(t)
	defer cleanup()

	mng := GetBackupManager(paths)

	timestamp := "20240928T1430"

	timestampMetaDir := filepath.Join(paths.Meta, timestamp)
	if err := os.MkdirAll(timestampMetaDir, 0755); err != nil {
		t.Fatalf("Failed to create timestamp meta dir: %v", err)
	}

	initialNote := ProtectionNote{
		Note:   "Initial backup",
		Reason: "Testing",
	}
	initialData, _ := yaml.Marshal(initialNote)
	infoPath := filepath.Join(timestampMetaDir, "info.yaml")
	if err := os.WriteFile(infoPath, initialData, 0644); err != nil {
		t.Fatalf("Failed to write initial info.yaml: %v", err)
	}

	backup := Backup{
		Timestamp:   timestamp,
		IsProtected: true,
		Items:       []Snapshot{},
	}

	err := mng.Restore(backup, nil)
	if err != nil {
		t.Errorf("Restore() unexpected error: %v", err)
	}

	data, err := os.ReadFile(infoPath)
	if err != nil {
		t.Fatalf("Failed to read info.yaml: %v", err)
	}

	var note ProtectionNote
	if err := yaml.Unmarshal(data, &note); err != nil {
		t.Fatalf("Failed to unmarshal info.yaml: %v", err)
	}

	if len(note.RestorationDates) != 1 {
		t.Errorf("Expected 1 restoration date, got %d", len(note.RestorationDates))
	}
}

func TestRestore_ExecutesPostRestoreHook(t *testing.T) {
	tempDir, paths, cleanup := setupRestoreTestDirs(t)
	defer cleanup()

	mng := GetBackupManager(paths)

	markerFile := filepath.Join(tempDir, "hook_executed")
	createHookScript(t, paths.Hooks, "post-restore", "touch "+markerFile)

	backup := Backup{
		Timestamp:   "20240928T1430",
		IsProtected: false,
		Items:       []Snapshot{},
	}

	err := mng.Restore(backup, nil)
	if err != nil {
		t.Errorf("Restore() unexpected error: %v", err)
	}

	if _, err := os.Stat(markerFile); os.IsNotExist(err) {
		t.Error("Post-restore hook was not executed")
	}
}

func TestRestore_ReturnsErrorOnHookFailure(t *testing.T) {
	_, paths, cleanup := setupRestoreTestDirs(t)
	defer cleanup()

	mng := GetBackupManager(paths)

	createHookScript(t, paths.Hooks, "post-restore", "exit 1")

	backup := Backup{
		Timestamp:   "20240928T1430",
		IsProtected: false,
		Items:       []Snapshot{},
	}

	err := mng.Restore(backup, nil)
	if err == nil {
		t.Error("Restore() expected error on hook failure, got nil")
	}
}
