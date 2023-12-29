package operations

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Mega-Kranus/NastySecrets/internal/consts"
	"github.com/Mega-Kranus/NastySecrets/internal/faults"
	"github.com/Mega-Kranus/NastySecrets/internal/validators"
)

type ConfigFile struct {
	Key          string            `json:"Key"`
	DoRename     bool              `json:"DoRename"`
	RenamedFiles map[string]string `json:"RenamedFiles"`
}

// This function retrieves the key from the config file
func retrieveKey(path string) (key []byte, err error) {
	// Check if the config file provided exists
	err = validators.ConfigExists()
	if err != nil {
		return nil, err
	}

	// Read the raw JSON data
	rawData, err := os.ReadFile(path)
	if err != nil {
		err = errors.Join(err, faults.GetError(consts.ECantOpenFile))
		return nil, err
	}

	// Deserialize raw JSON into a struct
	var keyFileDeserialized ConfigFile
	json.Unmarshal(rawData, &keyFileDeserialized)

	// Decode key from base64 to bytes
	keyInBytes, err := b64ToBytes(keyFileDeserialized.Key)
	if err != nil {
		return nil, err
	}

	// Check if the key is valid
	err = validators.IsKeyValid(&keyInBytes)
	if err != nil {
		return nil, err
	}

	return keyInBytes, nil
}

func b64ToBytes(b64Value string) (bytesValue []byte, err error) {
	bytesValue, err = base64.StdEncoding.DecodeString(b64Value)
	if err != nil {
		return nil, err
	}

	return bytesValue, nil
}

func bytesTob64(bytesValue []byte) (b64Value string) {
	return base64.StdEncoding.EncodeToString(bytesValue)
}

// This function will export all the data required for future decryption
func writeConfig(output string, key []byte, doRename bool, renamedMap map[string]string) (err error) {
	fmt.Println("[+] Writing to config...")
	// Initialize a 'configFile' object
	var dataToBeWritten ConfigFile

	// Assign the data to the struct object
	dataToBeWritten.Key = bytesTob64(key)
	dataToBeWritten.DoRename = doRename
	dataToBeWritten.RenamedFiles = renamedMap

	// Serialize the data into JSON format
	jsonData, err := json.Marshal(dataToBeWritten)
	if err != nil {
		return err
	}

	// Write the JSON data to the output path
	// Permission: -rw-r--r--
	err = os.WriteFile(output, jsonData, 0644)
	if err != nil {
		return err
	}

	fmt.Println("[+] Completed")
	fmt.Printf("\nThe config file was written to %v\n", output)
	fmt.Println("This file includes all the data necessary for decryption.")
	fmt.Println("It is recommended to keep it somewhere safe and offline.")

	return nil
}

// This function will import all the data required for decryption
func readConfig(configPath string, configData *ConfigFile) (err error) {
	fmt.Println("[+] Reading config...")

	// Read JSON data
	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	// Import data
	err = json.Unmarshal(jsonData, configData)
	if err != nil {
		return err
	}

	return nil
}
