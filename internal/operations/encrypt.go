package operations

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/faults"
	"github.com/Mega-Kranus/NastySecrets/internal/validators"
)

/*
Contains the paths of encrypted files
Currently used for emergency decryption if anything goes wrong
Package level variable
*/
var encryptedFiles []string

// This function generates a brand new 32 Bytes key
func generateNewKey(key []byte) (err error) {
	_, err = rand.Read(key)
	if err != nil {
		return err
	}

	return nil
}

// This function prepares for the encryption process
func InitiateEncryption(path, configPathKey string, doRename bool, threadsAmount int) (err error) {
	// Validate the output path
	configOutput, err := validators.ValidateOutputPath()
	if err != nil {
		return err
	}

	// Initialize key variable
	key := make([]byte, 32)

	// Check if a key is provided
	if len(configPathKey) > 0 {
		key, err = retrieveKey(configPathKey)
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
	filesAmount, err := gatherFiles(path)
	if err != nil {
		return err
	}

	// Return an error if no files are in the directory
	if filesAmount == 0 {
		err = faults.GetError(consts.ENoFilesFound)
		return err
	}

	// Write the key to the config file before encryption
	err = writeKey(configOutput, key)
	if err != nil {
		return err
	}

	// Execute the encryption process
	renamedMap, err := encrypt(key, foundFiles, filesAmount, threadsAmount, doRename, configOutput)
	if err != nil {
		return err
	}

	// Overwrite the config file with the key and the rest of the data
	err = writeConfig(configOutput, key, doRename, renamedMap)
	if err != nil {
		return err
	}

	return nil
}

// This function encrypts all files within the provided directory recursively and concurrently
func encrypt(key []byte, files []string, filesAmount, threadsAmount int, doRename bool, configOutput string) (renamedMap map[string]string, err error) {
	// Initializing a counter variable for going over the "files" slice
	var filesIterator int

	// Initialize the block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Initialize AES GCM
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Initialize a wait group
	var wg sync.WaitGroup

	/*
		Using a pointer to retrieve errors if they occure instead of channels
		Channels make the program run extremely slow (7x slower)
		Not sure if I'm doing something wrong but for now, I'm using a pointer
	*/
	var pError error

	for filesIterator < filesAmount {
		// Decide how many files to encrypt at a time
		// This helps avoid creating extra go routines if there are less than 8 files left
		atATime := threadsAmount
		if filesIterator+threadsAmount > filesAmount {
			atATime = filesAmount - filesIterator
		}

		wg.Add(atATime)
		for ; atATime > 0; atATime-- {
			go encryptConcurrent(&aesgcm, &wg, files[filesIterator], &pError)

			filesIterator++
		}
		wg.Wait()

		fmt.Printf("\r[+] Encrypted [%d/%d]", filesIterator, filesAmount)

		// See if at least one of the go routines had a problem
		// If so stop encrypting further and decrypt all encrypted files back
		if pError != nil {
			err = emergencyDecrypt(key, threadsAmount)
			if err != nil {
				// Join both errors and return them back
				pError = errors.Join(pError, err)
			}

			return nil, pError
		}
	}

	// Used for output formatting purposes
	fmt.Println()

	// This map will contain the renamed file name and the encrypted original name
	// example: temp123: {encryptedname.txt}
	renamedMap = make(map[string]string)

	if doRename {
		// Rename files
		err = encryptNames(&aesgcm, filesAmount, renamedMap)
		if err != nil {
			return nil, err
		}
	}

	return renamedMap, nil
}

// This function will encrypt a file and remove the go routine from its wait group
func encryptConcurrent(aesgcm *cipher.AEAD, wg *sync.WaitGroup, filePath string, pError *error) {
	defer (*wg).Done() // Remove from wait group

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

	// Add the path to the encrypted files slice
	encryptedFiles = append(encryptedFiles, filePath)
}
