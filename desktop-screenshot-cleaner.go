package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
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
			datePart := strings.Split(nameWithoutExt, " at ")[0]
			_, err := time.Parse("2006-01-02", strings.TrimPrefix(datePart, "Screenshot "))
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

var rootCmd = &cobra.Command{
	Use:   "screenshot-finder",
	Short: "Finds screenshots in a specified directory and copies them to another directory",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dirToSearch := args[0]
		destDir := args[1]

		err := os.MkdirAll(destDir, 0755)
		if err != nil {
			fmt.Println("error creating directory:", err)
			return
		}

		err = filepath.WalkDir(dirToSearch, listScreenshot(destDir))
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&deleteAfterCopy, "delete", false, "Delete screenshots after copying")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
