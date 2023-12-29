package operations

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/faults"
	"github.com/Mega-Kranus/NastySecrets/internal/validators"
)

// This function generates a brand new 32 Bytes key
func generateNewKey(key []byte) (err error) {
	_, err = rand.Read(key)
	if err != nil {
		return err
	}

	return nil
}

// This function encrypts all the files within the root directory recursively
func Encrypt(path, keypath string, doRename bool) (err error) {

	// Validate the output path
	configOutput, err := validators.ValidateOutputPath()
	if err != nil {
		return err
	}

	// Initialize key variable
	var key = make([]byte, 32)

	// Check if a key is provided
	if len(keypath) > 0 {
		key, err = retrieveKey(keypath)
		if err != nil {
			return err
		}
	} else {
		// Generate a new key
		err = generateNewKey(key)
		if err != nil {
			return err
		}
	}

	// Get a list of all files recursively
	foundFilesAmount, err := gatherFiles(path)
	if err != nil {
		return err
	}

	// Return an error if no files are in the directory
	if foundFilesAmount == 0 {
		err = faults.GetError(consts.ENoFilesFound)
		return err
	}

	// Encrypted files counter
	var encryptedFilesAmount int

	// Initialize the block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Initialize AES GCM
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Begin encryption
	for _, v := range FoundFiles {
		// Create a new random IV
		iv := make([]byte, aesgcm.NonceSize())
		if _, err = rand.Read(iv); err != nil {
			return err
		}

		// Read the file's data
		plainText, err := os.ReadFile(v)
		if err != nil {
			return err
		}

		// Encrypt the plain text and put the IV and the cipher text in a new slice
		ivAndCipherText := append(iv, aesgcm.Seal(nil, iv, plainText, nil)...)

		// Write the encrypted data to the file
		// Permission: -rw-r--r-- (Only if file doesn't exists and a new file must be created)
		err = os.WriteFile(v, ivAndCipherText, 0644)
		if err != nil {
			return err
		}

		encryptedFilesAmount++

		fmt.Printf("\r[+] Encrypted [%d/%d]", encryptedFilesAmount, foundFilesAmount)
	}

	// Used for output formatting purposes
	fmt.Println()

	var renamedMap = make(map[string]string)

	if doRename {
		// rename files
		err = encryptNames(&aesgcm, foundFilesAmount, renamedMap)
		if err != nil {
			return nil
		}
	}

	err = writeConfig(configOutput, key, doRename, renamedMap)
	if err != nil {
		return err
	}

	return nil
}
