package screenshot_cleaner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var deleteAfterCopy bool

func listScreenshot(destDir string) func(string, os.DirEntry, error) error {
	var counter int
	return func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Println("error:", err)
			return nil
		}
		if !d.IsDir() && (filepath.Ext(path) == ".png" || filepath.Ext(path) == ".jpg") && strings.Contains(strings.ToLower(d.Name()), "screenshot") {
			nameWithoutExt := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
			datePart := strings.Split(nameWithoutExt, " ")[1]
			_, err := time.Parse("2006-01-02", datePart)
			if err == nil {
				destPath := filepath.Join(destDir, d.Name())
				err := copyFile(path, destPath)
				if err != nil {
					fmt.Println("error copying file:", err)
				} else {
					counter++
					if deleteAfterCopy && !strings.HasPrefix(path, destDir) {
						err := os.Remove(path)
						if err != nil {
							fmt.Println("error deleting file:", err)
						}
					}
				}
			}
		}
		if counter == 0 {
			fmt.Println("No screenshot files found.")
		} else {
			fmt.Printf("Total number of screenshot files found: %d\n", counter)
		}
		return nil
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	return dstFile.Sync()
}
