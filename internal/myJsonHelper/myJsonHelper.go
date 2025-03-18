package myJsonHelper

import (
	"encoding/json"
	"os"
)

func ReadJSON(filename string) (map[string]interface{}, error) {
    // Read the file
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    // Parse JSON
    var result map[string]interface{}
    if err := json.Unmarshal(data, &result); err != nil {
        return nil, err
    }

    return result, nil
}