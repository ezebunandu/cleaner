package cleaner

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const (
	CmdUsage = iota
	CmdMove
)

type Cmd struct {
	Operation int
	Source    string
	Target    string
	Extension string
}

func NewCmd() *Cmd {
	return &Cmd{Extension: ""}
}

func CmdFromArgs(args []string) Cmd {
	if len(args) < 3 {
		return Cmd{Operation: CmdUsage}
	}
	fSet := flag.NewFlagSet("cleaner", flag.ContinueOnError)
	ext := fSet.String("ext", "", "the file extension to match")
	err := fSet.Parse(args[1:]) // skip program name so flags like -ext are parsed
	if err != nil {
		return Cmd{Operation: CmdUsage}
	}
	pos := fSet.Args()
	if len(pos) < 2 {
		return Cmd{Operation: CmdUsage}
	}
	return Cmd{Operation: CmdMove, Source: pos[0], Target: pos[1], Extension: *ext}
}

const usage = `usage: cleaner <SOURCE> <TARGET>`

// ListScreenshots lists screenshot files in dir.
func ListScreenshots(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var results = make([]string, 0, len(files))
	for _, file := range files {
		fname := file.Name()
		if strings.HasPrefix(fname, "Screenshot ") {
			results = append(results, dir+"/"+fname)
		}
	}
	return results, nil
}

// MoveScreenshot moves file to target.
func MoveScreenshot(file, target string) error {
	fileName := filepath.Base(file)
	dateSubfolder := DateSubfolder(fileName)
	targetPath := filepath.Join(target, dateSubfolder)

	_, err := os.Stat(targetPath)

	if err != nil {
		err := os.Mkdir(targetPath, 0700)
		if err != nil {
			return err
		}
	}

	targetName := filepath.Join(targetPath, fileName)
	err = os.Rename(file, targetName)
	if err != nil {
		return err
	}
	return nil
}

// DateSubfolder returns the date from filename.
func DateSubfolder(filename string) string {
	parts := strings.Split(filename, " ")

	if len(parts) < 2 {
		panic("DateSubfolder should only be called with file names matching a pattern like 'Screenshot 2024-07-30 at 9.55.08 AM.png'")
	}
	return parts[1]
}

func FileModeDateFromMetaData(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return info.ModTime().UTC().Format("2006-01-02"), nil
}

func ListFilesByExt(fsys fs.FS, ext string) ([]string, error) {
	if ext == "" {
		return nil, errors.New("extension cannot be empty")
	}
	// Normalise: strip leading dot, lowercase
	ext = strings.TrimPrefix(ext, ".")
	ext = strings.ToLower(ext)

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.ToLower(filepath.Ext(name)) == "."+ext {
			paths = append(paths, name)
		}
	}
	return paths, nil
}

func MoveFileByExt(file, target string) error {
	dateSubfolder, err := FileModeDateFromMetaData(file)
	if err != nil {
		return err
	}
	destDir := filepath.Join(target, dateSubfolder)
	if err := os.MkdirAll(destDir, 0700); err != nil {
		return err
	}
	destPath := filepath.Join(destDir, filepath.Base(file))
	if _, err := os.Stat(destPath); err == nil {
		fmt.Printf("skipping %s: already exists at %s\n", filepath.Base(file), destPath)
		return nil
	}
	return os.Rename(file, destPath)
}

func Main() int {
	cmd := CmdFromArgs(os.Args)
	if cmd.Operation == CmdUsage {
		fmt.Println(usage)
		return 0
	}

	if cmd.Extension != "" {
		// Ext-based path
		files, err := ListFilesByExt(os.DirFS(cmd.Source), cmd.Extension)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if len(files) == 0 {
			fmt.Println("no files to move")
			return 0
		}
		if _, err := os.Stat(cmd.Target); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		for _, f := range files {
			if err := MoveFileByExt(filepath.Join(cmd.Source, f), cmd.Target); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
		}
		fmt.Printf("moved %d files to %s\n", len(files), cmd.Target)
	} else {
		// Screenshot path (unchanged)
		screenshots, err := ListScreenshots(cmd.Source)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		if len(screenshots) == 0 {
			fmt.Println("no files to move")
			return 0
		}
		if _, err := os.Stat(cmd.Target); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
		for _, screenshot := range screenshots {
			if err := MoveScreenshot(screenshot, cmd.Target); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return 1
			}
		}
		fmt.Printf("moved %d files to %s\n", len(screenshots), cmd.Target)
	}
	return 0
}
