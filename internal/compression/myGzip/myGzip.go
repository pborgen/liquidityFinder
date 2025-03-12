package myGzip

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}
	gz.Close()
	return buf.Bytes(), nil
}

func DecompressJson(data []byte) ([]byte, error) {
	var result []byte

	 // Decompress the data
	 gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	 if err != nil {
		 return result, err
	 }
	 defer gzipReader.Close()

	 err = json.NewDecoder(gzipReader).Decode(&result)
	 if err != nil {
		 return result, err
	 }
	 
	 return result, nil
}
