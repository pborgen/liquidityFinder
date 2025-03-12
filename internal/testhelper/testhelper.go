package testhelper

import (
	"fmt"
	"os"
	"strings"
)

func Setup() {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}
	
	// find "goLangArb/"
	stringToFind := "goLangArb"
	stringToFindLength := len(stringToFind)
	index := strings.Index(dir, stringToFind)

	
	baseDir := dir[:index + stringToFindLength]

    os.Setenv("BASE_DIR", baseDir)
}
