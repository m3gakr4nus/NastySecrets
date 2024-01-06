package operations

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"os"
	"sync"

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

// This function encrypts all files within the provided directory recursively and concurrently
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

	// Initialize a wait group
	var wg sync.WaitGroup

	/*
		Using a pointer to retrieve errors if they occure instead of channels
		Channel make the the program run extremely slow (7x slower)
		Not sure if I'm doing something wrong but for now, i'm using a pointer
	*/
	var pError error

	for encryptedFilesAmount < foundFilesAmount {
		// Decide how many file to encrypt at a time
		// This helps avoid creating extra go routines if there are less than 8 files left
		atATime := 8
		if encryptedFilesAmount+atATime > foundFilesAmount {
			atATime = foundFilesAmount - encryptedFilesAmount
		}

		wg.Add(atATime)
		for ; atATime > 0; atATime-- {
			go encryptConcurrent(&aesgcm, &wg, FoundFiles[encryptedFilesAmount], &pError)

			encryptedFilesAmount++
		}
		wg.Wait()

		// See if at least one of the go routines had a problem, if so stop encrypting further
		if pError != nil {
			return pError
		}

		fmt.Printf("\r[+] Encrypted [%d/%d]", encryptedFilesAmount, foundFilesAmount)
	}

	// Used for output formatting purposes
	fmt.Println()

	var renamedMap = make(map[string]string)

	if doRename {
		// Rename files
		err = encryptNames(&aesgcm, foundFilesAmount, renamedMap)
		if err != nil {
			return err
		}
	}

	err = writeConfig(configOutput, key, doRename, renamedMap)
	if err != nil {
		return err
	}

	return nil
}

// This function will encrypt a file and remove the go routine from its wait group
func encryptConcurrent(aesgcm *cipher.AEAD, wg *sync.WaitGroup, filePath string, pError *error) {
	defer wg.Done() // remvoe from wait group

	// Create a new random IV
	iv := make([]byte, (*aesgcm).NonceSize())
	if _, err := rand.Read(iv); err != nil {
		*pError = err
		return
	}

	// Read the file's data
	plainText, err := os.ReadFile(filePath)
	if err != nil {
		*pError = err
		return
	}

	// Encrypt the plain text and put the IV and the cipher text in a new slice
	ivAndCipherText := append(iv, (*aesgcm).Seal(nil, iv, plainText, nil)...)

	// Write the encrypted data to the file
	// Permission: -rw-r--r-- (Only if file doesn't exists and a new file must be created)
	err = os.WriteFile(filePath, ivAndCipherText, 0644)
	if err != nil {
		*pError = err
		return
	}
}
