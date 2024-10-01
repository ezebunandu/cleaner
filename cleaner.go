package cleaner

import (
	"os"
	"strings"
	"path/filepath"
	"fmt"
)

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
	return strings.Split(filename, " ")[1]
}

func Main() int {
	if len(os.Args) != 3 {
        fmt.Println(usage)
        return 0
    }
    source, target := os.Args[1], os.Args[2]
    screenshots, err := ListScreenshots(source)
    if err != nil {
        fmt.Println(err)
        return 1
    }

    if len(screenshots) == 0 {
        fmt.Println("no files to move")
        return 0
    }

    _, err = os.Stat(target)

    if err != nil {
        fmt.Println(err)
        return 1
    }
    for _, screenshot := range screenshots{
        err := MoveScreenshot(screenshot, target)
        if err != nil {
           fmt.Println(err)
           return 1
        }
    }
    fmt.Printf("moved %d files to %s\n", len(screenshots), target)
	return 0
}