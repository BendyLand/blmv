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
	"time"
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
	fmt.Print("Calculating number of partitions: ")
	numPartitions := CalculateNumPartitions(len(files), 50, 2)
	time.Sleep(1 * time.Second)
	fmt.Println(numPartitions)
	currentPartition := 0
	fmt.Println("Dividing files...")
	time.Sleep(1 * time.Second)
	subDirs := partitionFiles(files, numPartitions)
	var wg sync.WaitGroup
	errChan := make(chan []os.File, numPartitions)
	fmt.Printf("Beginning move...\n\n")
	time.Sleep(1 * time.Second)
	for _, dir := range subDirs {
		wg.Add(1)
		go func(d []fs.DirEntry) {
			defer wg.Done()
			currentPartition++
			errFiles := beginMove(paths, d, currentPartition, numPartitions)
			errChan <- errFiles
		}(dir)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()
	var errFilesList []os.File
	for errs := range errChan {
		errFilesList = append(errFilesList, errs...)
	}
	if len(errFilesList) > 0 {
		fmt.Println("The following files may not have copied successfully:")
		for _, errFile := range errFilesList {
			fmt.Println(errFile)
		}
	} else {
		fmt.Println("All files moved successfully!")
	}
}

func beginMove(paths Paths, partition []fs.DirEntry, currentPartition int, numPartitions int) []os.File {
	fmt.Printf("Starting partition %d/%d.\n", currentPartition, numPartitions)
	src := paths.Src
	dst := paths.Dst
	count := 0
	var errorFiles []os.File
	for _, file := range partition {
		fmt.Printf("\nCurrent file: %s\n\n", file.Name())
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
			fmt.Println("Copied successfully: ", file.Name())
		}
		if !fileExists(dstFile.Name()) {
			fmt.Println("Possible error during copy. Adding file to error list...")
			errorFiles = append(errorFiles, *dstFile)
			continue
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
	fmt.Printf("\n%d files moved successfully! (Partition: %d/%d)\n\n", count, currentPartition, numPartitions)
	return errorFiles
}

func flatten(nested [][]os.File) []os.File {
	var result []os.File
	for _, innerSlice := range nested {
		result = append(result, innerSlice...)
	}
	return result
}

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	// An unexpected error occurred
	return false
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
