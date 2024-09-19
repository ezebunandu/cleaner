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

// func MoveScreenshot(file, target string) error {
// 	fileName := filepath.Base(file)
// 	targetName := filepath.Join(target, fileName)
// 	err := os.Rename(file, targetName)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

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