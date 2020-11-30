package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// Duplicate represents a duplicate file found in a directory.
type Duplicate struct {
	Files []File // An array of duplicate file names.
}

// Scanner provides methods for scanning a single directory.
type Scanner struct {
	directoryPath string
	maxThreads    uint
}

// NewScanner creates and returns a new Scanner struct with
// the provided directory path and maximum number of threads.
func NewScanner(directoryPath string, maxThreads uint) *Scanner {
	return &Scanner{
		directoryPath,
		maxThreads,
	}
}

// Scan spawns goroutines to scan the directory for duplicate files
// and returns an array of Duplicate structs with any duplicates found.
//
// Scan will block/pause until the scan is completed.
func (s *Scanner) Scan() ([]Duplicate, error) {
	// Create a wait group, which will allow us to wait on all our
	// scanning processes before returning a result.
	wg := &sync.WaitGroup{}

	// Create a new ScanMap.
	scanMap := NewScanMap()

	// Read all the file names of the supplied directory.
	files, err := ioutil.ReadDir(s.directoryPath)
	if err != nil {
		return nil, err
	}

	// Get the amount of files in the directory.
	numFiles := len(files)

	// If there are no files in the directory, return an
	// empty array because there are obviously no duplicates.
	if numFiles <= 0 {
		return []Duplicate{}, nil
	}

	// Set the number of goroutines to the specified maxThreads
	// parameter.
	numGoroutines := s.maxThreads
	// If we have less files than maximum threads, set the number of
	// goroutines equal to the number of files.
	if len(files) < int(s.maxThreads) {
		numGoroutines = uint(len(files))
	}

	// chunkSize is the number of files we want to process in each chunk.
	chunkSize := numFiles / int(numGoroutines)

	for i := 0; i < int(numGoroutines); i++ {
		// Add 1 to the wait group to wait for each process to
		// finish.
		wg.Add(1)

		// Get the starting index by multiplying the chunk size
		// times the index.
		start := i * chunkSize
		// Get the ending index by adding one chunk's size to the
		// start index.
		end := start + chunkSize

		// If we are on the last goroutine,
		// take all remaining files to ensure
		// none are left behind.
		if i == int(numGoroutines-1) {
			end = len(files)
		}

		chunk := files[start:end]

		// Spawn the chunk scanner as a goroutine to make it process
		// asynchronously.
		go s.scanChunk(chunk, scanMap, wg)
	}

	// Tell the function to wait until all the chunks are processed.
	wg.Wait()

	// Create a new array to push any found duplicates into.
	duplicates := []Duplicate{}
	// Iterate through the scan map and find any keys with more than one value
	// which are duplicates.
	sm := scanMap.GetMap()
	for _, files := range sm {
		// If one or less files are found, skip.
		if len(files) <= 1 {
			continue
		}

		// Create a Duplicate struct to append to our duplicates array.
		duplicate := Duplicate{
			Files: files,
		}

		// Append our duplicate to our duplicates.
		duplicates = append(duplicates, duplicate)
	}

	return duplicates, nil
}

// scanChunk takes an array of os.FileInfos and scans them.
func (s *Scanner) scanChunk(chunk []os.FileInfo, scanMap *ScanMap, wg *sync.WaitGroup) {
	for _, file := range chunk {
		// Ignore directories.
		if file.IsDir() {
			continue
		}

		// Get the full file path by joining the directory path
		// with the individual file's name.
		path := filepath.Join(s.directoryPath, file.Name())
		// Open the file for reading.
		f, err := os.Open(path)
		// If an error occurs, skip this file.
		if err != nil {
			fmt.Printf("error: couldn't open file %s\n", file.Name())
			continue
		}

		// Get the file's hash.
		hash, err := hashFile(f)
		// If an error occurs, skip this file.
		if err != nil {
			fmt.Printf("error: couldn't hash file %s\n", file.Name())
			continue
		}

		// Add the hash and the file name to the scan map.
		scanMap.Set(Hash(hash), File(file.Name()))
	}

	// Tell the wait group that we're done (subtracts
	// 1 from the wait group queue).
	wg.Done()
}
