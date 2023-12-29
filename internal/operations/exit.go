package operations

import (
	"fmt"
	"os"
)

// Print the error before exiting on fatal errors
func ExitOnError(err *error) {
	fmt.Println((*err).Error())
	os.Exit(1)
}
