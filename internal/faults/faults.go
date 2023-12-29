package faults

import (
	"errors"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
)

// Error codes and their coresponding error messages
var errorsAndMessages = map[uint]string{
	consts.EEncryptionAndDecryption:  "Error: Encryption and Decryption can not be used together!",
	consts.ENoEncryptionOrDecryption: "Error: Missing operation flag (use -e or -d)",
	consts.EConfigFileNotFound:       "Error: Unable to find the config file (-k {path_to_file})",
	consts.ECantOpenFile:             "Error: Could not open the config file!",
	consts.EInvalidKeyLength:         "Error: The provided key is invalid! Key length must be 16, 24 or 32 Bytes.",
	consts.ENoPathProvided:           "Error: Missing path (use -p {path_to_folder})",
	consts.EPathInvalid:              "Error: The provided path does not exist!",
	consts.EPathIsFile:               "Error: The provided path is a file. Need a directory.",
	consts.ENoFilesFound:             "Error: The path provided does not contain any files.",
	consts.EOutputPathInvalid:        "Error: The output path provided is invalid!",
}

// Returns the error message as an error
func GetError(errorCode uint) (err error) {
	errMessage, ok := errorsAndMessages[errorCode]
	if !ok {
		panic("The error code does not exists!")
	}

	err = errors.New(errMessage)

	return err
}
