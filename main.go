package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "screenshot-finder [source directory] [destination directory]",
	Short: "Finds screenshots in a specified directory and copies them to another directory",
	Long: `screenshot-finder is a CLI tool that finds screenshot files in the specified source directory 
and copies them to the specified destination directory.

Example usage:
screenshot-finder /path/to/source /path/to/destination --delete

This will find all screenshot files in /path/to/source, copy them to /path/to/destination, 
and delete the original files.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		dirToSearch := args[0]
		destDir := args[1]

		err := os.MkdirAll(destDir, 0755)
		if err != nil {
			fmt.Println("error creating directory:", err)
			return
		}

		err = filepath.WalkDir(dirToSearch, screenshot_cleaner.listScreenshot(destDir))
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
