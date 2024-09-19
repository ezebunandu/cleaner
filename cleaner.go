package cleaner

import (
	"os"
	"strings"
	"path/filepath"
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

func MoveScreenshot(file, target string) error {
	fileName := filepath.Base(file)
	targetName := filepath.Join(target, fileName)
	err := os.Rename(file, targetName)
	if err != nil {
		return err
	}
	return nil
}
