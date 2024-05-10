package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"strconv"
	"sync"
)

func main() {
	fmt.Println("Hello icmv!")
	srcPath := ""
	dstPath := ""
	contents := getDirContents(srcPath)
	splitContents := splitDir(contents, 100)
	beginConcurrentOperation(splitContents, srcPath, dstPath)
}

func beginConcurrentOperation(splitContents [][]fs.DirEntry, srcPath string, dstPath string) {
	var wg sync.WaitGroup
	for i, subDir := range splitContents {
		fmt.Printf("Writing to tempDir %d\n", i)
		wg.Add(1)
		go beginCopy(subDir, srcPath, dstPath, i, &wg)
	}
	wg.Wait()
	fmt.Println("Copy finished successfully!")
}

// THIS WILL NEED TO BE EDITED FOR TEMPORARY STORAGE CAPACITY
// OR TARGET THE SSD DIRECTLY
func beginCopy(subDir []fs.DirEntry, srcPath string, dstPath string, dirNum int, wg *sync.WaitGroup) error {
	defer wg.Done()
	num := strconv.Itoa(dirNum)
	err := os.Mkdir(dstPath+num, 0755)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return err
	}
	for _, item := range subDir {
		src, err := os.Open(srcPath + "/" + item.Name())
		if err != nil {
			fmt.Printf("Error opening %v: %v\n", item.Name(), err)
			return err
		}
		defer src.Close()

		dst, err := os.Create(dstPath + num + "/" + item.Name())
		if err != nil {
			fmt.Printf("Error creating destination: %v\n", err)
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		if err != nil {
			fmt.Printf("Error copying from src to dst: %v\n", err)
			return err
		}
	}
	return nil
}

func calculateStartingIndices(length int, pieceLength int) []int {
	var result []int
	for i := 0; i < length; i += pieceLength {
		result = append(result, i)
	}
	return result
}

func splitDir(contents []fs.DirEntry, pieceLength int) [][]fs.DirEntry {
	var result [][]fs.DirEntry
	length := len(contents)
	for start := 0; start < length; start += pieceLength {
		end := start + pieceLength
		if end > length {
			end = length
		}
		result = append(result, contents[start:end])
	}
	return result
}

func getDirContents(path string) []fs.DirEntry {
	contents, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Unable to read directory: %v\n", err)
		os.Exit(1)
	}
	return contents
}
