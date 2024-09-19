package cleaner_test

// get list of files in the directory
// get the files we are interested in
// copy files

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/ezebunandu/cleaner"
)

func TestListScreenshots_CorrectlyListsScreenshotsinDirectory(t *testing.T) {
	t.Parallel()
	want := []string{"testdata/Screenshot 2024-07-30 at 9.55.08 AM.png"}

	got, err := cleaner.ListScreenshots("testdata")
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(want, got) {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestListFiles_ReturnsErrorWhenDirNotReadable(t *testing.T) {
	t.Parallel()
	_, err := cleaner.ListScreenshots("bogus")
	if err == nil {
		t.Error("want error when directory unreadable, got nil")
	}
}

func TestMoveScreenshot_CopiesScreenshotToTargetDir(t *testing.T) {
	t.Parallel()
	target := t.TempDir()
	file := t.TempDir() + "/" + "Screenshot 2024-07-30 at 9.55.08 AM.png"
	want := []byte{1, 2, 3}
	err := os.WriteFile(file, want, 0o600)
	if err != nil {
		t.Fatal(err)
	}	
	err = cleaner.MoveScreenshot(file, target)
	if err != nil {
		t.Fatal(err)
	}
	destPath := target + "/" + "Screenshot 2024-07-30 at 9.55.08 AM.png"
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Fatalf("expected file at %s but it does not exist", destPath)
	}
	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatal(err)
	}
	if ! slices.Equal(want, got) {
		t.Error("target does not contain the data in source")
	}
}

func TestMoveScreenshot_RemovesScreenshotFromSourcDir(t *testing.T){
	t.Parallel()
	target := t.TempDir()
	source := t.TempDir()
	screenshotFile := "Screenshot 2024-07-30 at 9.55.08 AM.png"
	sourcePath := source + "/" + screenshotFile
	file, err := os.Create(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	cleaner.MoveScreenshot(sourcePath, target)
	if err!= nil {
		t.Fatal(err)
	}
	want := []string{}
	got, err := cleaner.ListScreenshots(source)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(want, got) {
		t.Errorf("want %q got %q", want, got)
	}
}

func TestMoveScreenshot_CreatesTargetDirIfNotAlreadyExisting(t *testing.T){
	t.Parallel()
	source := t.TempDir()

	screenshotFile := "Screenshot 2024-07-30 at 9.55.08 AM.png"
	sourcePath := filepath.Join(source, screenshotFile)
	target := filepath.Join(t.TempDir(), "bogus")

	err := cleaner.MoveScreenshot(sourcePath, target)
	if err != nil {
		t.Fatalf("MoveScreenshot returned an error: %v", err)
	}
	if _, err := os.Stat(target); os.IsNotExist(err) {
		t.Fatalf("expected target directory %q to be created, but it doesn't exist", target)
	}
}