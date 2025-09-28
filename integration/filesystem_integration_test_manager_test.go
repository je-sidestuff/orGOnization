package integration

import (
	"os"
	"testing"
)

func TestFilesystemIntegrationTestManager(t *testing.T) {
	t.Parallel()

	// Create a new instance of FilesystemTestInput
	fti := NewFilesystemTestInput()

	// Copy the test fixture directory to a new temporary directory
	_, err := fti.CloneDirectoryTreeToNewTempFolder("../test/fixtures/file_trees/direct/simple", "test_fixture")

	if err != nil {
		t.Fatalf("Failed to clone directory tree: %v", err)
	}

	// Delete the temporary directory to clean up when we are done
	defer RunTestStage(t, "cleanup", func() { fti.CleanupFilesystemTestInput() }, func() {
		t.Logf("Skipping cleanup of directory %s", fti.GetTempFolderPath("test_fixture"))
	})

	// Verify that the temporary directory exists and contains the expected files with the correct content
	expectedFilesToContentMap := map[string]string{
		"root_file.txt":              "Root file content.",
		"nested_dir/nested_file.txt": "Nested file content.",
	}
	for file, content := range expectedFilesToContentMap {
		_, err := os.Stat(fti.GetTempFolderPath("test_fixture") + "/" + file)
		if err != nil {
			t.Errorf("Failed to stat file %s: %v", file, err)
		}

		contentFromTempdir, err := os.ReadFile(fti.GetTempFolderPath("test_fixture") + "/" + file)
		if err != nil {
			t.Errorf("Failed to read file %s: %v", file, err)
		}

		if string(contentFromTempdir) != content {
			t.Errorf("Expected file %s to have content %q, but got %q", file, expectedFilesToContentMap[file], string(contentFromTempdir))
		}
	}
}
