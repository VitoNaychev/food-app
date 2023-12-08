package testutil

import "fmt"

func HandleLoadEnviornmentError(err error) {
	fmt.Printf("couln't load enviornment, got error: %v\n", err)
}
