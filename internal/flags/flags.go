package flags

import "flag"

// All command line arguments will be stored in these variables
var (
	FlagEncrypt bool
	FlagDecrypt bool
	FlagRename  bool
	FlagConfig  string
	FlagOutput  string
	FlagPath    string
)

// This is the initialization function to set the flags and parse them
func SetFlags() {
	flag.BoolVar(&FlagEncrypt, "e", false, "Perform encryption")
	flag.BoolVar(&FlagDecrypt, "d", false, "Perform decryption")
	flag.BoolVar(&FlagRename, "n", false, "Rename files")
	flag.StringVar(&FlagConfig, "k", "", "Config file's path for decryption or to re-use an old key")
	flag.StringVar(&FlagOutput, "o", "./", "Config file's output path")
	flag.StringVar(&FlagPath, "p", "", "The folder to encrypt/decrypt its data recursively")

	flag.Parse()
}
