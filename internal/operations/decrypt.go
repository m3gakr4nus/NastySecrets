package operations

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"os"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/faults"
)

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

	// Beging decryption
	for _, v := range FoundFiles {
		// Read encrypted file's data
		cipherText, err := os.ReadFile(v)
		if err != nil {
			return err
		}

		// Extract the IV from the data
		iv := cipherText[:aesgcm.NonceSize()]
		cipherText = cipherText[aesgcm.NonceSize():]

		// Decrypt the data
		decryptedData, err := aesgcm.Open(nil, iv, cipherText, nil)
		if err != nil {
			return err
		}

		// Write the decrypted data back
		// Permission: -rw-r--r-- (Only if file doesn't exists and a new file must be created)
		err = os.WriteFile(v, decryptedData, 0644)
		if err != nil {
			return err
		}

		decryptedFilesAmount++

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
