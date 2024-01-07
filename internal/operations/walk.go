package operations

import (
	"fmt"
	"os"
	"path/filepath"
)

/*
All found file paths will be stored here after walking the directory
Package level variable
*/
var foundFiles []string

// This function will add each file to the slice
func visit(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		foundFiles = append(foundFiles, path)
		return nil
	}

	return nil
}

// This function will walk the provided directory
func gatherFiles(path string) (int, error) {
	fmt.Printf("[*] Gathering files...\n")

	err := filepath.Walk(path, visit)
	if err != nil {
		return 0, err
	}

	fmt.Printf("[+] Found %d\n", len(foundFiles))

	return len(foundFiles), nil
}
