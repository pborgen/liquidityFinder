package myEncrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	"io"
	"os"

	"github.com/pborgen/liquidityFinder/internal/myConfig"
	"github.com/pborgen/liquidityFinder/internal/myUtil"
)


const encryptedKeyFileName = "encryptedKey.txt"


func Setup() {


	// Encode the byte slice to a hexadecimal string
	randomKey := generateRandomKey()

	seedPw := myConfig.GetInstance().SeedPw
	
	encryptedKeyString, err := encrypt(randomKey, []byte(seedPw))
	if err != nil {
		panic(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// Store the seed in a file
	file, err := os.Create(homeDir + "/" + encryptedKeyFileName)
	if err != nil {
		panic(err)
	}
	defer file.Close() // Ensure the file is closed when done

	// Write content to the file
	_, err = file.WriteString(encryptedKeyString)
	if err != nil {
		panic(err)
	}
}

func Encrypt(plainText string) (string, error) {
	return encrypt([]byte(plainText), []byte(getKey()))
}

func Decrypt(cipherTextHex string) (string, error) {
	return decrypt(cipherTextHex, []byte(getKey()))
}

func encrypt(plaintext, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func decrypt(ciphertextHex string, key []byte) (string, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		panic("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err)	
	}

	return string(plaintext), nil
}

func getKey() string {

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	encryptedSeedString, err := myUtil.ReadFileToString(homeDir + "/" + encryptedKeyFileName)
	if err != nil {
		panic(err)
	}

	seedPw := myConfig.GetInstance().SeedPw

	seedString, err := decrypt(encryptedSeedString, []byte(seedPw))
	if err != nil {
		panic(err)
	}

	return seedString
}

func generateRandomKey() ([]byte) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}
