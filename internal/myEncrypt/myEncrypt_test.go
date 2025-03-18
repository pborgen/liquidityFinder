package myEncrypt

import (
	"testing"

	"github.com/paulborgen/goLangArb/internal/testhelper"
)
func TestSetup(t *testing.T) {

	testhelper.Setup()
	Setup()
}


func TestEncrypt(t *testing.T) {

	testhelper.Setup()
	key := "4af318424513a82c0245e5458ff82101701f6f9d05320d87a3bde8d"

	encrypted, err := Encrypt(key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	decrypted, err := Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decrypted != key {
		t.Fatalf("Decrypted text does not match original text")
	}

	
}
