package validators

import (
	"os"
	"path/filepath"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/faults"
	"github.com/Mega-Kranus/NastySecrets/internal/flags"
)

// Check if the path provided exists and is valid
func ValidatePath() (err error) {
	// Check if path is provided
	if len(flags.FlagPath) > 0 {
		// Retrieve path info
		info, err := os.Stat(flags.FlagPath)

		// Check if path exists
		if os.IsNotExist(err) {
			err = faults.GetError(consts.EPathInvalid)
			return err
		} else if err != nil {
			return err
		}

		// Check if path is a directory
		if !info.IsDir() {
			err = faults.GetError(consts.EPathIsFile)
			return err
		}

		return nil
	}

	// If no path was provided
	err = faults.GetError(consts.ENoPathProvided)
	return err
}

// This function validates the output path provided for the config file
func ValidateOutputPath() (outputPath string, err error) {
	// TODO: Check user permissions as well
	if len(flags.FlagOutput) > 2 {
		pathInfo, err := os.Stat(flags.FlagOutput)
		if os.IsNotExist(err) {
			// If provided path doesn't exist
			_, err := os.Stat(filepath.Dir(flags.FlagOutput))
			if err != nil {
				return "", faults.GetError(consts.EOutputPathInvalid)
			}

			// Path provided contains a to-be-create file
			return flags.FlagOutput, nil
		} else if err != nil {
			return "", err
		}

		// If user provided a directory as output
		if pathInfo.IsDir() {
			// Join directory and filename
			outputPath = filepath.Join(flags.FlagOutput, "nasty")
		} else {
			// User provided an already existing file as output
			outputPath = flags.FlagOutput
		}
	} else {
		// If no path was provided, default will be current directory
		outputPath = "./nasty"
	}

	return outputPath, nil
}

// Identify wether to encrypt or decrypt
func IdentifyOperation() (operation uint, err error) {
	switch {
	case flags.FlagEncrypt && !flags.FlagDecrypt:
		// If only encryption flag is set (valid)
		operation = consts.Encryption
		err = nil
	case flags.FlagDecrypt && !flags.FlagEncrypt:
		// If only decryption flag is set (valid)
		operation = consts.Decryption
		err = nil
	case flags.FlagEncrypt && flags.FlagDecrypt:
		// If both flags are set (invalid)
		operation = 0
		err = faults.GetError(consts.EEncryptionAndDecryption)
	case !flags.FlagEncrypt && !flags.FlagDecrypt:
		// If neither of the flags are set (invalid)
		operation = 0
		err = faults.GetError(consts.ENoEncryptionOrDecryption)
	}

	return operation, err
}

// Check if config file exists
func ConfigExists() (err error) {
	// Check if file can be opened/exists
	f, err := os.OpenFile(flags.FlagConfig, os.O_RDWR, 0644)
	if err != nil {
		err = faults.GetError(consts.EConfigFileNotFound)
	}
	defer f.Close()

	return err
}

// Check if key is valid
func IsKeyValid(keyInBytes *[]byte) (err error) {
	// Get the key length
	keyLength := len(*keyInBytes)

	// Compare key length with valid key lengths
	for _, v := range consts.ValidKeyLengths {
		if keyLength == v {
			// Key valid
			return nil
		}
	}

	// Key invalid
	err = faults.GetError(consts.EInvalidKeyLength)

	return err
}

// Make sure about renaming while decryption
// The user can forget to provided the flag (-n). The files will then be decrypted without
