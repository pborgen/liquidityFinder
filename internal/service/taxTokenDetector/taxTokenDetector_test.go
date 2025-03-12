package taxTokenDetector

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestTaxTokenDetector(t *testing.T) {

    dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

    // Find goLangArb directory
    goLangArbDir := filepath.Join(dir, "..", "..", "..")

    os.Setenv("BASE_DIR", goLangArbDir)

   log.Println("Current working directory:", dir)
    
    td, err := NewTaxDetector()
    if err != nil {
        t.Fatalf("Failed to create tax detector: %v", err)
    }


    x, z := td.DetectTaxToken(
        context.Background(), 
        common.HexToAddress("0xeE2208C3552b573A895f2E5CcF08A4D010F951F3"))

    fmt.Println(x, z)

}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
