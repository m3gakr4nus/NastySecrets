package operations

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"sync"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/faults"
)

// This function decrypts all files within the provided directory recursively and concurrently
func Decrypt(path, configPath string) (err error) {
	var configData ConfigFile

	// Validate if config file is provided
	if !(len(configPath) > 0) {
		err = faults.GetError(consts.EConfigFileNotFound)
		return err
	}

	// Read config files data
	readConfig(configPath, &configData)

	// Retrieve key
	key, err := b64ToBytes(configData.Key)
	if err != nil {
		return err
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

	// Decrypted files counter
	var decryptedFilesAmount int

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

	// Beging decryption
	for decryptedFilesAmount < foundFilesAmount {
		atATime := 8
		if decryptedFilesAmount+atATime > foundFilesAmount {
			atATime = foundFilesAmount - decryptedFilesAmount
		}

		wg.Add(atATime)
		for ; atATime > 0; atATime-- {
			go decryptConcurrent(&aesgcm, &wg, FoundFiles[decryptedFilesAmount], &pError)

			decryptedFilesAmount++
		}
		wg.Wait()

		if pError != nil {
			return pError
		}

		fmt.Printf("\r[+] Decrypted [%d/%d]", decryptedFilesAmount, foundFilesAmount)
	}

	// Used for output formatting purposes
	fmt.Println()

	// Rename files back if necessary
	if configData.DoRename {
		err = decryptNames(&aesgcm, foundFilesAmount, configData.RenamedFiles)
		if err != nil {
			return err
		}
	}

	fmt.Println("[+] Completed")

	return nil
}

// This function will decrypt a file and remove the go routine from its wait group
func decryptConcurrent(aesgcm *cipher.AEAD, wg *sync.WaitGroup, filePath string, pError *error) {
	defer wg.Done()

	// Read encrypted file's data
	cipherText, err := os.ReadFile(filePath)
	if err != nil {
		*pError = err
		return
	}

	// Extract the IV from the data
	iv := cipherText[:(*aesgcm).NonceSize()]
	cipherText = cipherText[(*aesgcm).NonceSize():]

	// Decrypt the data
	decryptedData, err := (*aesgcm).Open(nil, iv, cipherText, nil)
	if err != nil {
		*pError = err
		return
	}

	// Write the decrypted data back
	// Permission: -rw-r--r-- (Only if file doesn't exists and a new file must be created)
	err = os.WriteFile(filePath, decryptedData, 0644)
	if err != nil {
		*pError = err
		return
	}
}
