package main

import (
	"crypto/sha1"
	"io"
	"os"
)

// hashFile takes in a pointer to an os.File and hashes the file's
// contents, returning a string representation of the hash sum.
func hashFile(file *os.File) (string, error) {
	// Initializes a new SHA-1 hasher.
	hash := sha1.New()

	// Make a buffer so we can sequentially load in from our
	// file.
	buf := make([]byte, 10*1024)
	for {
		// Loop and take at most 10*1024 bytes of data
		// from the file.
		numBytes, err := file.Read(buf)
		if numBytes > 0 {
			// If we read more than 0 bytes, aka we have
			// new data, write it into the hash.
			_, err := hash.Write(buf[:numBytes])
			if err != nil {
				return "", err
			}
		}

		// Check if the reader has reached the end of the file.
		// If so, break out of the loop so we can return our
		// sum.
		if err == io.EOF {
			break
		}

		// If the error is not an end-of-file error,
		// exit the function and return.
		if err != nil {
			return "", err
		}
	}

	// Return the final hash sum.
	return string(hash.Sum(nil)), nil
}
