package main

import (
	"flag"
	"fmt"
	"os"
)

// Our command line arguments are listed here.
// In "advanced" applications, one might opt to use
// a struct rather than global variables, or to
// not use global variables alltogether.

// directoryName sets the name of the directory that
// will be scanned. It is treated as a relative path.
var directoryName string

// threadCount sets the maximum number of "threads" or
// goroutines to spawn during scanning.
var threadCount uint

// init runs before the program enters the main function.
// It is not used much outside of program setup, such as
// command line flags.
func init() {
	flag.StringVar(&directoryName, "path", "", "a path to a directory to scan")
	flag.UintVar(&threadCount, "threads", 10, "the maximum number of threads to use during scanning")
}

// main is the entry point of our program, and runs all
// the main logic of it.
func main() {
	// Parse our flags.
	flag.Parse()

	// If no directory name is provided, return an error.
	if directoryName == "" {
		fmt.Println("error: please pass a valid directory name")
		os.Exit(1)
	}

	// Create a new Scanner with the given directory and thread count.
	scanner := NewScanner(directoryName, threadCount)
	// Run the scanner (this function will block until the scanner finishes).
	dupes, err := scanner.Scan()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print a nice finished message.
	fmt.Println("- finished --------------")

	// If no duplicates were found, tell the user and exit.
	if len(dupes) <= 0 {
		fmt.Println("  no duplicate files found!")
		fmt.Println("  have a nice day :)")
		return
	}

	// Otherwise, print how many duplications were found.
	fmt.Printf("-> %d duplications found.\n\n", len(dupes))

	// Iterate through each of the duplicates.
	for index, dupe := range dupes {
		// Print the duplicate index in the duplicates array.
		fmt.Printf("• group %d\n", index)

		// Iterate through all the duplicate files.
		for _, file := range dupe.Files {
			// Print the file name.
			fmt.Printf("  |- %s\n", file)
		}

		// Print a newline to add some space below each group.
		fmt.Print("\n")
	}
}
