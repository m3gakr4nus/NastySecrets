package operations

import (
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

// This function will rename all the files to temp with a counter (temp1, temp2, temp3 and so on)
func encryptNames(aesgcm *cipher.AEAD, foundFilesAmount int, renamedMap map[string]string) (err error) {
	var renamedCounter int

	// Beging encrypting names
	for _, v := range foundFiles {
		renamedCounter++

		// Get the file's basename
		fileName := []byte(filepath.Base(v))
		fileDir := filepath.Dir(v)

		// Initialize the new name and path
		newName := fmt.Sprintf("temp%d", renamedCounter)
		newPath := filepath.Join(fileDir, newName)

		// Generate a new random IV
		iv := make([]byte, (*aesgcm).NonceSize())
		if _, err := rand.Read(iv); err != nil {
			return err
		}

		// Encrypt the file's name and put the IV and the encrypted name in a slice
		ivAndCipherName := append(iv, (*aesgcm).Seal(nil, iv, fileName, nil)...)

		// Add the old name and the corresponding encrypted name (as base64) to the map
		renamedMap[newName] = base64.StdEncoding.EncodeToString(ivAndCipherName)

		// Rename the file
		err = os.Rename(v, newPath)
		if err != nil {
			return err
		}

		fmt.Printf("\r[+] Renamed [%d/%d]", renamedCounter, foundFilesAmount)
	}

	// Output formatting purposes
	fmt.Println()

	return nil
}

func decryptNames(aesgcm *cipher.AEAD, foundFilesAmount int, renamedMap map[string]string) (err error) {
	var renamedCounter int

	// Beging decrypting names
	for _, v := range foundFiles {
		renamedCounter++

		// Get the file's basename
		fileName := filepath.Base(v)
		fileDir := filepath.Dir(v)

		// Retrieve the cipher name
		cipherName, ok := renamedMap[string(fileName)]
		if !ok {
			fmt.Printf("Warning: File (%s) could not be renamed! Skipping...\n", v)
			continue
		}

		// Delete the name after retrieval
		delete(renamedMap, fileName)

		// Decode the cipher name
		plainCipherName, err := base64.StdEncoding.DecodeString(cipherName)
		if err != nil {
			fmt.Printf("Warning: File (%s) could not be decoded. Skipping...\n", v)
			continue
		}

		// Extract IV from the cipher name
		iv := plainCipherName[:(*aesgcm).NonceSize()]
		plainCipherName = plainCipherName[(*aesgcm).NonceSize():]

		// Decrypt the name
		oldName, err := (*aesgcm).Open(nil, iv, plainCipherName, nil)
		if err != nil {
			return err
		}

		// Initialize the new path for renaming
		newPath := filepath.Join(fileDir, string(oldName))

		// Rename the file
		err = os.Rename(v, newPath)
		if err != nil {
			return err
		}

		fmt.Printf("\r[+] Renamed [%d/%d]", renamedCounter, foundFilesAmount)
	}

	// Output formatting purposes
	fmt.Println()

	return nil
}
