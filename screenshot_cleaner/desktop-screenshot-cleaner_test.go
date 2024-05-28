package screenshot_cleaner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListScreenshot(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create some test files in the temporary directory
	testFiles := []string{
		"Screenshot 2022-01-01 at 10.00.00.png",
		"Screenshot 2022-01-02 at 12.30.45.jpg",
		"not_a_screenshot.png",
		"Screenshot_without_date.png",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tempDir, file)
		_, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Call the listScreenshot function with the temporary directory
	listFunc := listScreenshot(tempDir)

	// Call the listFunc with each test file
	err := filepath.WalkDir(tempDir, listFunc)
	if err != nil {
		t.Fatalf("listScreenshot failed: %v", err)
	}

	// TODO: Add assertions to verify the expected behavior of listScreenshot
	// For example, you can check if the expected files were copied to the destination directory.
}
