package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	paths, err := os.ReadFile("config.txt")
	if err != nil {
		fmt.Println("There was a problem getting the paths.")
		return
	}
	lines := strings.Split(string(paths), "\n")
	directory := lines[0] 
	destination := lines[1]
	moveDownloadedFiles(directory, destination)
}

func moveDownloadedFiles(srcPath string, dstPath string) error {
	contents, err := os.ReadDir(srcPath)
	if err != nil {
		fmt.Println("Error getting contents")
		return err
	}
	fmt.Println(contents)
	return nil
}