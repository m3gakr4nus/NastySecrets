package validators

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/faults"
	"github.com/Mega-Kranus/NastySecrets/internal/flags"
)

// Validate if all flags that are needed by all operations are correct
func ValidateGlobalFlags() (err error) {
	// Check if a path is provided and is valid
	err = ValidatePath()
	if err != nil {
		return err
	}

	// Check if the threads amount is more than 0
	err = ValidateThreadsAmount()
	if err != nil {
		return err
	}

	return nil
}

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
	/*
		TODO: 	Add posibility for users to export config file in a
				1 character long file name
				currently if (-o 'o') output will be --> ./nasty
				output file name gets ignore if it's less than 2 chars
	*/
	if flags.FlagOutput != "./" {
		pathInfo, err := os.Stat(flags.FlagOutput)
		if os.IsNotExist(err) {
			// If provided path doesn't exist check parent directory's content
			_, err := os.Stat(filepath.Dir(flags.FlagOutput))
			if err != nil {
				// If also the parent directory doesn't exist
				return "", faults.GetError(consts.EOutputPathInvalid)
			}

			// Path provided contains a to-be-create file
			return flags.FlagOutput, nil
		} else if err != nil {
			// This typically handles permission errors and other
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

// Identify if an action flag has been provided (-e or -d)
// If so what is it (encryption or decryption)
// If not provided, return an error
func IdentifyOperation() (operation uint, err error) {
	switch {
	case flag.NFlag() < 1:
		// If no flags are provided (invalid)
		operation = 0
		err = faults.GetError(consts.ENoFlagProvided)
	case flags.FlagVersion && flag.NFlag() <= 1:
		// If only -v is provided (valid)
		operation = consts.ShowVersion
		err = nil
	case flags.FlagVersion && flag.NFlag() > 1:
		// If the -v flag is provided with other flags (invalid)
		operation = 0
		err = faults.GetError(consts.EVersionFlagNotAlone)
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

// Check if the threads amount is between 0 and 26 (1-25)
func ValidateThreadsAmount() (err error) {
	switch {
	case flags.FlagThreads < 1:
		return faults.GetError(consts.EThreadsAmountLessThanOne)
	case flags.FlagThreads > 25:
		return faults.GetError(consts.EThreadsAmountUnsafe)
	default:
		return nil
	}
}
