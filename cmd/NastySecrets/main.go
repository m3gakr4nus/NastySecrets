package main

import (
	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/flags"
	"github.com/Mega-Kranus/NastySecrets/internal/operations"
	"github.com/Mega-Kranus/NastySecrets/internal/validators"
)

func main() {
	// Initialize and set CLI flags
	flags.SetFlags()

	operation, err := validators.IdentifyOperation()
	if err != nil {
		// Can not continue without knowing the operation to perform
		operations.ExitOnError(&err)
	}

	// Check if a path is provided and is valid
	err = validators.ValidatePath()
	if err != nil {
		operations.ExitOnError(&err)
	}

	// Perform operations accordingly
	switch operation {
	case consts.Encryption:
		// Beging encryption
		err = operations.InitiateEncryption(flags.FlagPath, flags.FlagConfig, flags.FlagRename)
		if err != nil {
			operations.ExitOnError(&err)
		}
	case consts.Decryption:
		// Beging decryption
		err = operations.InitiateDecryption(flags.FlagPath, flags.FlagConfig)
		if err != nil {
			operations.ExitOnError(&err)
		}
	}
}
