package main

import (
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
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
	numPartitions := CalculateNumPartitions(len(files), 50, 2)
	subDirs := partitionFiles(files, numPartitions)
	var wg sync.WaitGroup
	for _, dir := range subDirs {
		wg.Add(1)
		go func(d []fs.DirEntry) {
			defer wg.Done()
			beginMove(paths, d, &wg)
		}(dir)
	}
	wg.Wait()
}

func beginMove(paths Paths, partition []fs.DirEntry, wg *sync.WaitGroup) {
	fmt.Println("Beginning move...")
	src := paths.Src
	dst := paths.Dst
	count := 0
	for _, file := range partition {
		fmt.Println("Current file:", file.Name())
		srcPath := filepath.Join(src, file.Name())
		dstPath := filepath.Join(dst, file.Name())

		srcFile, err := os.Open(srcPath)
		if err != nil {
			fmt.Println("Problem opening file:", err)
			continue
		}

		dstFile, err := os.Create(dstPath)
		if err != nil {
			fmt.Println("Problem creating destination:", err)
			continue
		}

		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()
		if err != nil {
			fmt.Println("Problem copying file:", err)
			continue
		} else {
			fmt.Println("Copied successfully:", file.Name())
		}

		err = os.Remove(srcFile.Name())
		if err != nil {
			fmt.Println("Problem removing source file:", err)
			continue
		} else {
			fmt.Println("Removed successfully:", file.Name())
		}

		count += 1
	}
	fmt.Printf("%d files moved successfully!\n", count)
	wg.Done()
}

func partitionFiles(files []fs.DirEntry, numPartitions int) [][]fs.DirEntry {
	partitionSize := (len(files) + numPartitions - 1) / numPartitions
	partitions := make([][]fs.DirEntry, 0, numPartitions)
	for i := 0; i < len(files); i += partitionSize {
		end := i + partitionSize
		if end > len(files) {
			end = len(files)
		}
		partitions = append(partitions, files[i:end])
	}
	return partitions
}

// Calculates the number of partitions to make based on the # of files.
// k is a constant representing a reasonable batch size per partition, and
// c is a tuning constant to adjust the growth rate (e.g., 2 or 4, depending on how aggressively you want to increase partitions).
func CalculateNumPartitions(numFiles int, k float64, c float64) int {
	numPartitions := math.Max(1, math.Min(float64(numFiles)/k, c*math.Sqrt(float64(numFiles))))
	return int(numPartitions)
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
