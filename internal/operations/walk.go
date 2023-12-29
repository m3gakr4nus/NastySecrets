package operations

import (
	"fmt"
	"os"
	"path/filepath"
)

// Found files will be stored here
var FoundFiles []string

// This function will add each file to the slice
func visit(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		FoundFiles = append(FoundFiles, path)
		return nil
	}

	return nil
}

// This function will walk the root directory
func gatherFiles(path string) (int, error) {
	fmt.Printf("[+] Gathering files...\n")

	err := filepath.Walk(path, visit)
	if err != nil {
		return 0, err
	}

	fmt.Printf("[+] Found %v\n", len(FoundFiles))

	return len(FoundFiles), nil
}
