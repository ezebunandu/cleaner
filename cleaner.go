package cleaner

import (
	"os"
	"path/filepath"
	"strings"
)

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

func MoveScreenshots(screenshots []string, targetDir string) {
	// if screenshots is emtpy, then do nothing.
	// if targetDir does not exist, create it first but don't return an error
	// if any of the screenshots files are already in the target, then do nothing -- don't return an error
	for _, screenshot := range screenshots {
		fileName := filepath.Base(screenshot)
		targetPath := targetDir + "/" + fileName
		
		err := os.Rename(screenshot, targetPath)
		if err!= nil {
			continue
		}
	}
}
