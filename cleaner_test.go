package cleaner_test

import (
	"os"
	"slices"
	"testing"
	"fmt"

	"github.com/ezebunandu/cleaner"
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
	err := os.WriteFile(sourcePath, []byte{}, 0o600)
	if err != nil {
		t.Fatal(err)
	}
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

func TestDateSubfolder_ReturnsCorrectSubfolderGivenFileName(t *testing.T){
	t.Parallel()
	filename := "Screenshot 2024-07-30 at 9.55.08 AM.png"
	want := "2024-07-30"
	got := cleaner.DateSubfolder(filename)
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func ExampleListScreenshots(){
		
	got, err := cleaner.ListScreenshots("testdata")
	if err != nil {
		panic(err)
	}
	fmt.Println(got)
	// Output: 
	// [testdata/Screenshot 2024-07-30 at 9.55.08AM.png]

}