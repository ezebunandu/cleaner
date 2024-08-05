package cleaner_test

// get list of files in the directory
// get the files we are interested in
// copy files

import (
	"os"
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
	source := t.TempDir()
	screenshotFile := "Screenshot 2024-07-30 at 9.55.08 AM.png"
	sourcePath := source + "/" + screenshotFile
	s, err := cleaner.ListScreenshots(target) 
	if err != nil {
		t.Fatal(err)
	}
	if len(s)!= 0 {
		for _, v := range s {
			if v == screenshotFile {
				t.Fatal("screenshot file already in target")
			}
		}
	}
	err = cleaner.MoveScreenshots([]string{sourcePath}, target)
	if err!= nil {
		t.Fatal(err)
	}
	destPath := target + "/" + screenshotFile
	want := []string{destPath}
	got, _ := cleaner.ListScreenshots(target)
	if!slices.Equal(want, got) {
		t.Errorf("want %q, got %q", want, got)
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
	err = cleaner.MoveScreenshots([]string{sourcePath}, target)
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