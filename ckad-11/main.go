package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	sunCipherID := os.Getenv("SUN_CIPHER_ID")
	for {
		fmt.Printf("SUN_CIPHER_ID: %s\n", sunCipherID)
		time.Sleep(1 * time.Second)

	}
}
