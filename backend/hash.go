package main

import (
	"fmt"

	"github.com/aryasatyawa/bayarin/internal/pkg/crypto"
)

func main() {
	hash, _ := crypto.HashPassword("passwordceo")
	fmt.Println(hash)
}
