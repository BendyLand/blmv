package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func isDownloaded(file string) bool {
    cmd := exec.Command("stat", "-f", "%b", file)
    output, err := cmd.Output()
    if err != nil {
        return false
    }
    allocatedBlocks := strings.TrimSpace(string(output))
    return allocatedBlocks != "0" // checks for allocated bytes
}

func downloadFile(file string) error {
    cmd := exec.Command("brctl", "download", file)
    var out bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &out
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error downloading:", out.String())
    }
    return err
}

func walkPath(dir string, iCloudBasePath string) error {
    // Verify if the directory is within iCloud Drive
    if !strings.HasPrefix(dir, iCloudBasePath) {
        message := fmt.Sprintf("Error: The specified directory is not within iCloud Drive (%s).\n", iCloudBasePath)
        return errors.New(message)
    }

    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if !info.IsDir() {
            if !isDownloaded(path) {
                fmt.Printf("Downloading %s...\n", path)
                if err := downloadFile(path); err != nil {
                    fmt.Printf("Failed to download %s: %v\n", path, err)
                }
            } else {
                fmt.Printf("%s is already downloaded.\n", path)
            }
        }
        return nil
    })

	var message error
    if err != nil {
        temp := fmt.Sprintf("Error walking the path %s: %v\n", dir, err)
		message = errors.New(temp)
    }
	return message // either error message or nil
}

func main() {
	paths, err := os.ReadFile("config.txt")
	if err != nil {
		fmt.Println("There was a problem getting the paths.")
		return
	}
	lines := strings.Split(string(paths), "\n")
    iCloudBasePath := lines[0]
	directory := lines[1]
	walkPath(directory, iCloudBasePath)
}
