package consts

// NastySecret's information
const (
	Version   = "NastySecrets v1.3.0"
	UsageMenu = Version + `
	
Usage: nastysecrets [arguments] [options]
	
Arguments:
	-e		Perform encryption
	-d		Perform decryption
	-p		Folder to encrypt/decrypt its files recursively
	-v		Print the current version
Options:
	-n		Rename files to 'temp' for added privacy
			Default: false

	-t		How many files to encrpyt/decrypt simultaneously
			Default: 8
			Values: [1-25]
	
	-o		Path to write the config file to
			Default: . (Current directory)
	
	-c		Config file's path for decryption or re-using an old key
			Mandatory for decryption
	
More information: man nastysecrets
`
)

// Error enums
// Errors beging with 'E' for better accessability
const (
	EEncryptionAndDecryption = iota + 1
	ENoEncryptionOrDecryption
	EConfigFileNotFound
	ECantOpenFile
	EInvalidKeyLength
	ENoPathProvided
	EPathInvalid
	EPathIsFile
	ENoFilesFound
	EOutputPathInvalid
	EThreadsAmountLessThanOne
	EThreadsAmountUnsafe
	EVersionFlagNotAlone
	ENoFlagProvided
)

// Operation enum
const (
	Encryption = iota + 1
	Decryption
	ShowVersion
)

// Valid key lengths
var ValidKeyLengths = [3]int{16, 24, 32}
