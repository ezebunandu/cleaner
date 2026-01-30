package cleaner_test

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/ezebunandu/cleaner"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestListScreenshots_CorrectlyListsScreenshotsinDirectory(t *testing.T) {
	t.Parallel()
	want := []string{"testdata/Screenshot 2024-07-30 at 9.55.08AM.png"}

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
	destPath := target + "/2024-07-30/Screenshot 2024-07-30 at 9.55.08 AM.png"
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Fatalf("expected file at %s but it does not exist", destPath)
	}
	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(want, got) {
		t.Error("target does not contain the data in source")
	}
}

func TestMoveScreenshot_RemovesScreenshotFromSourcDir(t *testing.T) {
	t.Parallel()
	target := t.TempDir()
	source := t.TempDir()
	screenshotFile := "Screenshot 2024-07-30 at 9.55.08 AM.png"
	sourcePath := source + "/" + screenshotFile
	err := os.WriteFile(sourcePath, []byte{}, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	cleaner.MoveScreenshot(sourcePath, target)
	if err != nil {
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

func TestDateSubfolder_ReturnsCorrectSubfolderGivenFileName(t *testing.T) {
	t.Parallel()
	filename := "Screenshot 2024-07-30 at 9.55.08 AM.png"
	want := "2024-07-30"
	got := cleaner.DateSubfolder(filename)
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestFileModeDate_ReturnsCorrectFileModDateFromMetaData(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.txt")

	// Create an empty file
	_, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}

	// Set a fixed modification time
	fixedTime := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	err = os.Chtimes(path, fixedTime, fixedTime)
	if err != nil {
		t.Fatalf("failed to set file time: %v", err)
	}

	// Call the function
	got, err := cleaner.FileModeDateFromMetaData(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify the result in yyyy-mm-dd format
	want := "2024-01-15"
	if got != want {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestFileModeDate_ErrorsWhenFileNonExistent(t *testing.T) {
	t.Parallel()
	
	_, err := cleaner.FileModeDateFromMetaData("/nonexistent/file/path.txt")
	if err == nil {
		t.Error("want error for non-existent file, got nil")
	}
	
}
func ExampleListScreenshots() {

	got, err := cleaner.ListScreenshots("testdata")
	if err != nil {
		panic(err)
	}
	fmt.Println(got)
	// Output:
	// [testdata/Screenshot 2024-07-30 at 9.55.08AM.png]

}

func Test(t *testing.T) {
	t.Parallel()
	testscript.Run(t, testscript.Params{
		Dir: "testdata/testscript",
	})
}

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"cleaner": cleaner.Main,
	}))
}
