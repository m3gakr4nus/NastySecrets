package operations

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"
	"sync"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/faults"
	"github.com/Mega-Kranus/NastySecrets/internal/validators"
)

// This function prepares for the decryption process
func InitiateDecryption(path, configPath string, threadsAmount int) (err error) {
	var configData ConfigFile

	// Validate if config file is provided
	if len(configPath) <= 0 {
		// No config file was provided
		err = faults.GetError(consts.EConfigFileNotFound)
		return err
	}

	err = validators.ConfigExists()
	if err != nil {
		err = faults.GetError(consts.EConfigFileNotFound)
		return err
	}

	// Read config files data
	readConfig(configPath, &configData)

	// Decode the key from base64 to bytes
	key, err := b64ToBytes(configData.Key)
	if err != nil {
		return err
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

	// Execute the decryption process
	aesgcm, err := decrypt(key, foundFiles, filesAmount, threadsAmount)
	if err != nil {
		return err
	}

	// Rename files back if necessary
	if configData.DoRename {
		err = decryptNames(aesgcm, filesAmount, configData.RenamedFiles)
		if err != nil {
			return err
		}
	}

	return nil
}

// This function decrypts all files within the provided directory recursively and concurrently
func decrypt(key []byte, files []string, filesAmount, threadsAmount int) (*cipher.AEAD, error) {
	/*
		Decrypted files counter
		Used for going over the "files" slice
	*/
	var decryptedFilesAmount int

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

	// Beging decryption
	for decryptedFilesAmount < filesAmount {
		// Decide how many files to decrypt at a time
		// This helps avoid creating extra go routines if there are less than 8 files left
		atATime := threadsAmount
		if decryptedFilesAmount+threadsAmount > filesAmount {
			atATime = filesAmount - decryptedFilesAmount
		}

		wg.Add(atATime)
		for ; atATime > 0; atATime-- {
			go decryptConcurrent(&aesgcm, &wg, files[decryptedFilesAmount], &pError)

			decryptedFilesAmount++
		}
		wg.Wait()

		fmt.Printf("\r[+] Decrypted [%d/%d]", decryptedFilesAmount, filesAmount)

		if pError != nil {
			return nil, pError
		}
	}

	// Used for output formatting purposes
	fmt.Println()
	fmt.Println("[+] Completed")

	return &aesgcm, nil
}

// This function will decrypt a file and remove the go routine from its wait group
func decryptConcurrent(aesgcm *cipher.AEAD, wg *sync.WaitGroup, filePath string, pError *error) {
	defer (*wg).Done() // Remove from wait group

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

/*
This function decrypts the files that just got encrypted back
It gets executed if an error happens while encrypting a file
This helps to avoid further encryption and data loss.
Once an error is detected, it is unsafe to continue.
User must resolve the issue and try again
*/
func emergencyDecrypt(key []byte, threadsAmount int) (err error) {
	fmt.Println()
	fmt.Println("[!] An error occured during encryption!")

	// The amount of files that were encrypted (package level variable: "encryptedFiles")
	filesAmount := len(encryptedFiles)

	// If any files were encrypted, decrypt them back
	if filesAmount > 0 {
		fmt.Println("[!] Decrypting files back as it is unsafe to continue")

		_, err = decrypt(key, encryptedFiles, filesAmount, threadsAmount)
		if err != nil {
			return err
		}
	}

	return nil
}
