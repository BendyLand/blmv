package main

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
)

func main() {
	paths, err := getPaths()
	if err != nil {
		fmt.Println(err)
		return
	}
	files, err := os.ReadDir(paths.Src)
	if err != nil {
		fmt.Println(err)
		return
	}
	subDirs := partitionFiles(files, 50)
	for _, dir := range subDirs {
		fmt.Println(dir)
	}
}

func partitionFiles(files []fs.DirEntry, numPartitions int) [][]fs.DirEntry {
	partitionSize := len(files) / numPartitions
	if len(files) % partitionSize != 0 {
		partitionSize += 1
	}
	result := make([][]fs.DirEntry, numPartitions)
	for _, file := range files {
		temp := make([]fs.DirEntry, partitionSize)
		for i := range partitionSize {
			temp[i] = file
		}
		result = append(result, temp)
	}
	return result
}

func calculateNumPartitions(numFiles int) int {
	numPartitions := numFiles / 100
	if numPartitions < 1 {
		numPartitions = 1
	}
	return numPartitions
}

type Paths struct {
	Src string
	Dst string
}

func getPaths() (Paths, error) {
	paths, err := os.ReadFile("config.txt")
	if err != nil {
		message := fmt.Errorf("There was a problem getting the paths.\n")
		return Paths{}, message
	}
	lines := strings.Split(string(paths), "\n")
	source := lines[0]
	destination := lines[1]
	result := Paths{source, destination}
	return result, nil
}
