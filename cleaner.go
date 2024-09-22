package cleaner

import (
	"os"
	"strings"
	"path/filepath"
	"fmt"
)

const usage = `usage: cleaner <SOURCE> <TARGET>`

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

func DateSubfolder(filename string) string {
	return strings.Split(filename, " ")[1]
}

func Main(){
	if len(os.Args) != 3 {
        fmt.Println(usage)
        os.Exit(0)
    }
    source, target := os.Args[1], os.Args[2]
    screenshots, err := ListScreenshots(source)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if len(screenshots) == 0 {
        fmt.Println("no files to move")
        os.Exit(0)
    }

    _, err = os.Stat(target)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    for _, screenshot := range screenshots{
        err := MoveScreenshot(screenshot, target)
        if err != nil {
           fmt.Println(err)
           os.Exit(1)
        }
    }
    fmt.Printf("moved %d files to %s\n", len(screenshots), target)
}