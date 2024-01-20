package main

import (
	"fmt"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/flags"
	"github.com/Mega-Kranus/NastySecrets/internal/operations"
	"github.com/Mega-Kranus/NastySecrets/internal/validators"
)

func main() {
	// Initialize, set and parse CLI flags
	flags.SetFlags()

	// Validate if an action is provided and identify what it is
	operation, err := validators.IdentifyOperation()
	if err != nil {
		// Can not continue without knowing the operation to perform
		operations.ExitOnError(&err)
	}

	// Perform operations accordingly
	switch operation {
	case consts.Encryption:
		// Validate necessary flags
		err = validators.ValidateGlobalFlags()
		if err != nil {
			operations.ExitOnError(&err)
		}

		// Beging encryption
		err = operations.InitiateEncryption(flags.FlagPath, flags.FlagConfig, flags.FlagRename, flags.FlagThreads)
		if err != nil {
			operations.ExitOnError(&err)
		}
	case consts.Decryption:
		// Validate necessary flags
		err = validators.ValidateGlobalFlags()
		if err != nil {
			operations.ExitOnError(&err)
		}

		// Beging decryption
		err = operations.InitiateDecryption(flags.FlagPath, flags.FlagConfig, flags.FlagThreads)
		if err != nil {
			operations.ExitOnError(&err)
		}
	case consts.ShowVersion:
		// Show version
		fmt.Println(consts.Version)
	}
}
