package main

import "sync"

// Hash represents a hashsum of a file as a byte array.
type Hash string

// File represents a file's relative path as a string.
type File string

// ScanMap contains a hash map of hashes to file paths
// as well as a mutex for concurrent operation on the
// hash map.
type ScanMap struct {
	hashMap map[Hash][]File
	mutex   sync.Mutex
}

// NewScanMap creates and returns a new ScanMap by initializing
// a blank hash map and creating a new mutex.
func NewScanMap() *ScanMap {
	return &ScanMap{
		hashMap: map[Hash][]File{},
	}
}

// Set locks the internal mutex and writes the hash and file
// to the internal hash map.
func (sm *ScanMap) Set(hash Hash, file File) error {
	// Request access to the mutex, which will block until
	// access is acquired.
	sm.mutex.Lock()
	// We defer the unlock, which means it will be called
	// when the function exits, after the return statement
	// in order to unlock the mutex when we finish
	// modifying the hash map.
	defer sm.mutex.Unlock()

	// If there isn't a previously existing record for the hash,
	// intialize a new file array.
	if _, found := sm.hashMap[hash]; !found {
		sm.hashMap[hash] = []File{}
	}

	// Set the hash to the file in the hash map.
	sm.hashMap[hash] = append(sm.hashMap[hash], file)

	// Return nil as there is no error.
	return nil
}

// Length gives the current length of the map.
func (sm *ScanMap) Length() int {
	return len(sm.hashMap)
}

// GetMap returns the internal hash map.
func (sm *ScanMap) GetMap() map[Hash][]File {
	// Request access to the mutex, which will block until
	// access is acquired.
	sm.mutex.Lock()
	// We defer the unlock, which means it will be called
	// when the function exits, after the return statement
	// in order to unlock the mutex when we finish
	// modifying the hash map.
	defer sm.mutex.Unlock()

	return sm.hashMap
}
