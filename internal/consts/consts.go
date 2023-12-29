package consts

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
)

// Operation enum
const (
	Encryption = iota + 1
	Decryption
)

// Valid key lengths
var ValidKeyLengths = [3]int{16, 24, 32}
