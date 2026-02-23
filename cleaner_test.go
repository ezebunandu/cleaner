package cleaner_test

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"testing/fstest"
	"time"

	"github.com/ezebunandu/cleaner"
	"github.com/google/go-cmp/cmp"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestCmdFromArgs(t *testing.T) {
	testCases := []struct {
		desc string
		args []string
		want cleaner.Cmd
	}{
		{
			desc: "With No Args Returns Usage",
			args: []string{"cleaner"},
			want: cleaner.Cmd{Operation: cleaner.CmdUsage,},
		},
		{
			desc: "With Only Source Returns Usage",
			args: []string{"cleaner", "source"},
			want: cleaner.Cmd{Operation: cleaner.CmdUsage,},
		},
		{
			desc: "With Source and Target Returns Move",
			args: []string{"cleaner", "source", "target"},
			want: cleaner.Cmd{Operation: cleaner.CmdMove, Source: "source", Target: "target", Extension: ""},
		},
		{
			desc: "With Ext Flag Returns Move with Ext",
			args: []string{"cleaner", "-ext", "png", "source", "target"},
			want: cleaner.Cmd{Operation: cleaner.CmdMove, Source: "source", Target: "target", Extension: "png"},
		},	
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			got := cleaner.CmdFromArgs(tC.args)
			if !cmp.Equal(tC.want, got){
				t.Error(cmp.Diff(tC.want, got))
			}
		})
	}
}

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

func TestFilesByExt_CorrectlyListsFilesMatchingGivenExt(t *testing.T) {
	t.Parallel()
	fsys := fstest.MapFS{
		"file.png":                {},
		"subfolder/subfolder.png": {},
		"subfolder2/another.go":   {},
		"subfolder2/file.png":     {},
	}
	want := []string{"file.png"}
	got, err := cleaner.ListFilesByExt(fsys, "png")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestListFilesByExt_TableDriven(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		ext   string
		fsys  fstest.MapFS
		want  []string
	}{
		{
			name: "Top-level only",
			ext:  "png",
			fsys: fstest.MapFS{
				"file.png":     {},
				"sub/file.png": {},
			},
			want: []string{"file.png"},
		},
		{
			name: "Case-insensitive PNG uppercase",
			ext:  "png",
			fsys: fstest.MapFS{
				"photo.PNG": {},
			},
			want: []string{"photo.PNG"},
		},
		{
			name: "Case-insensitive mixed case",
			ext:  "png",
			fsys: fstest.MapFS{
				"img.Png": {},
			},
			want: []string{"img.Png"},
		},
		{
			name: "Dot normalisation",
			ext:  ".png",
			fsys: fstest.MapFS{
				"file.png": {},
			},
			want: []string{"file.png"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cleaner.ListFilesByExt(tt.fsys, tt.ext)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(tt.want, got) {
				t.Error(cmp.Diff(tt.want, got))
			}
		})
	}
}

func TestFilesByExt_ErrorsWhenExtNotGiven(t *testing.T) {
	t.Parallel()
	fsys := fstest.MapFS{
		"file.png": {},
	}
	_, err := cleaner.ListFilesByExt(fsys, "")
	if err == nil {
		t.Fatal("want error given empty ext by got none")
	}

}
func TestMoveFileByExt_MovesToCorrectDateSubfolder(t *testing.T) {
	t.Parallel()
	source := t.TempDir()
	target := t.TempDir()

	filePath := filepath.Join(source, "photo.png")
	if err := os.WriteFile(filePath, []byte{1, 2, 3}, 0o600); err != nil {
		t.Fatal(err)
	}

	fixedTime := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(filePath, fixedTime, fixedTime); err != nil {
		t.Fatal(err)
	}

	if err := cleaner.MoveFileByExt(filePath, target); err != nil {
		t.Fatal(err)
	}

	want := filepath.Join(target, "2024-03-15", "photo.png")
	if _, err := os.Stat(want); os.IsNotExist(err) {
		t.Fatalf("expected file at %s but it does not exist", want)
	}
}

func TestMoveFileByExt_SkipsAndWarnWhenDestinationExists(t *testing.T) {
	t.Parallel()
	source := t.TempDir()
	target := t.TempDir()

	filePath := filepath.Join(source, "photo.png")
	if err := os.WriteFile(filePath, []byte{1, 2, 3}, 0o600); err != nil {
		t.Fatal(err)
	}

	fixedTime := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(filePath, fixedTime, fixedTime); err != nil {
		t.Fatal(err)
	}

	// Pre-create destination with different content
	destDir := filepath.Join(target, "2024-03-15")
	if err := os.MkdirAll(destDir, 0700); err != nil {
		t.Fatal(err)
	}
	destPath := filepath.Join(destDir, "photo.png")
	if err := os.WriteFile(destPath, []byte{9, 9, 9}, 0o600); err != nil {
		t.Fatal(err)
	}

	if err := cleaner.MoveFileByExt(filePath, target); err != nil {
		t.Fatal(err)
	}

	// Source still exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("expected source file to still exist after skip, but it does not")
	}

	// Destination content unchanged
	got, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal([]byte{9, 9, 9}, got) {
		t.Error("destination content was modified unexpectedly")
	}
}

func TestMoveFileByExt_RemovesSourceAfterSuccessfulMove(t *testing.T) {
	t.Parallel()
	source := t.TempDir()
	target := t.TempDir()

	filePath := filepath.Join(source, "photo.png")
	if err := os.WriteFile(filePath, []byte{1, 2, 3}, 0o600); err != nil {
		t.Fatal(err)
	}

	if err := cleaner.MoveFileByExt(filePath, target); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("expected source file to be removed after move, but it still exists")
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
