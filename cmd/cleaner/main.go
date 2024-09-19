package main

import (
	"fmt"
	"os"

	"github.com/ezebunandu/cleaner"
)

const usage = `usage: cleaner <SOURCE> <TARGET>`

func main(){
    // get source and target from arguments
    // call ListScreenshots on the source to get slice of all screenshot files in source
    // range over the source and move each screenshot file to target

    if len(os.Args) != 3 {
        fmt.Println(usage)
        os.Exit(0)
    }
    source, target := os.Args[1], os.Args[2]
    screenshots, err := cleaner.ListScreenshots(source)
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
        err := cleaner.MoveScreenshot(screenshot, target)
        if err != nil {
           fmt.Println(err)
           os.Exit(1)
        }
    }
    fmt.Printf("moved %d files to %s\n", len(screenshots), target)
}