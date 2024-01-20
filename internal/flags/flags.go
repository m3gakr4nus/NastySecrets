package flags

import (
	"flag"
	"fmt"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
)

// All command line arguments will be stored in these variables
var (
	FlagEncrypt bool
	FlagDecrypt bool
	FlagRename  bool
	FlagVersion bool
	FlagConfig  string
	FlagOutput  string
	FlagPath    string
	FlagThreads int
)

// This is the initialization function to set the flags and parse them
func SetFlags() {
	flag.BoolVar(&FlagEncrypt, "e", false, "Perform encryption")
	flag.BoolVar(&FlagDecrypt, "d", false, "Perform decryption")
	flag.BoolVar(&FlagRename, "n", false, "Rename files")
	flag.BoolVar(&FlagVersion, "v", false, "Print the current version")
	flag.StringVar(&FlagConfig, "c", "", "Config file's path for decryption or to re-use an old key")
	flag.StringVar(&FlagOutput, "o", "./", "Config file's output path")
	flag.StringVar(&FlagPath, "p", "", "The folder to encrypt/decrypt its data recursively")
	flag.IntVar(&FlagThreads, "t", 8, "The number of files to encrypt/decrypt at a time")

	flag.Usage = flagUsage
	flag.Parse()
}

// This function prints out the customized usage menu (-h)
func flagUsage() {
	fmt.Printf("%s", consts.UsageMenu)
}
