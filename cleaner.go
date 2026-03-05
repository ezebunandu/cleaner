package cleaner

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
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
	MatchedFiles []string
	DryRun bool
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
	dryrun := fSet.Bool("dry-run", false, "enable dry-run mode")
	err := fSet.Parse(args[1:]) // skip program name so flags like -ext are parsed
	if err != nil {
		return Cmd{Operation: CmdUsage}
	}
	pos := fSet.Args()
	if len(pos) < 2 {
		return Cmd{Operation: CmdUsage}
	}
	return Cmd{Operation: CmdMove, Source: pos[0], Target: pos[1], Extension: *ext, DryRun: *dryrun}
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

// MoveFiles moves file to the datesubfolder on target.
func MoveFiles(file, target string) error {
	fileName := filepath.Base(file)
	dateSubfolder, err := DateSubfolder(fileName)
	targetPath := filepath.Join(target, dateSubfolder)

	_, err = os.Stat(targetPath)

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

var date = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)

// DateSubfolder returns the date from filename if the filename contains a date.
// else it gets the date from the file metadata's last modified time
func DateSubfolder(filename string) (string, error) {
	matches := date.FindStringSubmatch(filename)

	if len(matches) < 2 {
		return FileModeDateFromMetaData(filename)
	}
	return matches[1], nil
}

func FileModeDateFromMetaData(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	return info.ModTime().Format("2006-01-02"), nil
}

func ListFilesByExt(root string, fsys fs.FS, ext string, ) (paths []string, err error) {
	if ext == "" {
		return nil, errors.New("extension cannot be empty")
	}
	matchExt := ext
	if ext[0] != '.' {
		matchExt = "." + ext
	}
	fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if filepath.Ext(p) == matchExt {
			paths = append(paths, root+"/"+p)
		}
		return nil
	})
	return paths, nil
}

func printDryRun(out io.Writer, cmd Cmd) {
	fmt.Fprintf(out, "would have moved the following files from %s to %s\n", cmd.Source, cmd.Target)
	for _, f := range cmd.MatchedFiles {
		fmt.Fprintln(out, f)
	}
}

func matchedFilesForCmd(cmd Cmd) ([]string, error) {
	if cmd.Extension == "" {
		return ListScreenshots(cmd.Source)
	}
	return ListFilesByExt(cmd.Source, os.DirFS(cmd.Source), cmd.Extension)
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	cmd := CmdFromArgs(args)
	if cmd.Operation == CmdUsage {
		fmt.Fprintln(out, usage)
		return 0
	}
	files, err := matchedFilesForCmd(cmd)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	if len(files) == 0 {
		fmt.Fprintln(out, "no files to move")
		return 0
	}
	cmd.MatchedFiles = files
	if cmd.DryRun {
		printDryRun(out, cmd)
		return 0
	}
	_, err = os.Stat(cmd.Target)
	if err != nil {
		fmt.Fprintln(errOut, err)
		return 1
	}
	for _, f := range cmd.MatchedFiles {
		err := MoveFiles(f, cmd.Target)
		if err != nil {
			fmt.Fprintln(errOut, err)
			return 1
		}
	}
	fmt.Fprintf(out, "moved %d files to %s\n", len(cmd.MatchedFiles), cmd.Target)
	return 0
}

func Main() int {
	return run(os.Args, os.Stdout, os.Stderr)
}