package main

import (
	"fmt"
)

func main() {
	fmt.Printf("Starting Tests.\n")

	// add concurrent test
	// AddConcurrent()

	// synchronous test
	SyncCommunication()

	// asynchronous test
	AsyncCommunication()

	fmt.Printf("Completed Tests.\n")
}
