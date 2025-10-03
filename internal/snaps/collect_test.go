package snaps

import (
	"os"
	"path/filepath"
	"testing"
)

func createSnapshots(t *testing.T, baseDir string, snapshots []string) {
	t.Helper()
	for _, snap := range snapshots {
		snapPath := filepath.Join(baseDir, snap)
		err := os.Mkdir(snapPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create snapshot dir %s: %v", snap, err)
		}
	}
}

func createTestDir(t *testing.T, baseDir, name string) string {
	t.Helper()
	dirPath := filepath.Join(baseDir, name)
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test dir %s: %v", name, err)
	}
	return dirPath
}

func assertCollectResults(t *testing.T, result CollectResult, err error, expectedCount int, expectedGroupCount int) {
	t.Helper()
	if err != nil {
		t.Errorf("Collect() error = %v, want nil", err)
	}
	if result.TotalCount != expectedCount {
		t.Errorf("Collect() count = %d, want %d", result.TotalCount, expectedCount)
	}
	if expectedGroupCount == 0 && result.Groups != nil {
		t.Errorf("Collect() groups = %v, want nil", result.Groups)
		return
	}
	if len(result.Groups) != expectedGroupCount {
		t.Errorf("Collect() groups length = %d, want %d", len(result.Groups), expectedGroupCount)
	}
}

func TestCollect(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "btrbk_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	t.Run("empty directory", func(t *testing.T) {
		emptyDir := createTestDir(t, tempDir, "empty")
		manager := GetManagerForDirectory(emptyDir)
		result, err := manager.Collect()
		assertCollectResults(t, result, err, 0, 0)
	})

	t.Run("directory with valid snapshots", func(t *testing.T) {
		snapDir := createTestDir(t, tempDir, "snapshots")

		snapshots := []string{
			"@home.20240928T1430",
			"@var.20240928T1430",
			"@home.20240927T1400",
			"@root.20240928T1430",
		}

		createSnapshots(t, snapDir, snapshots)

		manager := GetManagerForDirectory(snapDir)
		result, err := manager.Collect()
		assertCollectResults(t, result, err, 4, 2)
	})

	t.Run("directory with mixed entries", func(t *testing.T) {
		mixedDir := createTestDir(t, tempDir, "mixed")

		validSnaps := []string{
			"@home.20240928T1430",
			"@var.20240927T1200",
		}
		createSnapshots(t, mixedDir, validSnaps)

		filePath := filepath.Join(mixedDir, "regular_file.txt")
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		file.Close()

		invalidDir := filepath.Join(mixedDir, "invalid_dir")
		err = os.Mkdir(invalidDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create invalid dir: %v", err)
		}

		manager := GetManagerForDirectory(mixedDir)
		result, err := manager.Collect()
		assertCollectResults(t, result, err, 4, 2)
	})

	t.Run("non-existent directory", func(t *testing.T) {
		nonExistentDir := filepath.Join(tempDir, "does_not_exist")
		manager := GetManagerForDirectory(nonExistentDir)
		result, err := manager.Collect()
		assertCollectResults(t, result, err, 0, 0)
	})
}
