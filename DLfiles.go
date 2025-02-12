/*DLfiles.go*/

package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DLFile represents a file with a name and SHA-256 hash.
type DLFile struct {
	Name string `json:"name"`
	Hash string `json:"hash"`
}

// hashFile calculates the SHA-256 hash of a single file.
func hashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashing file: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// hashDirectory calculates and returns the SHA-256 hash of every file in a directory.
func hashDirectory(dirPath string) (map[string]DLFile, error) {
	dlfiles := make(map[string]DLFile)

	// Walk the directory to process files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("accessing %s: %w", path, err)
		}

		// Skip directories and non-regular files
		if !info.Mode().IsRegular() {
			return nil
		}

		// Hash the file
		hash, err := hashFile(path)
		if err != nil {
			return fmt.Errorf("hashing %s: %w", path, err)
		}

		// Store the file info in the map
		dlfiles[info.Name()] = DLFile{
			Name: info.Name(),
			Hash: hash,
		}

		fmt.Printf("Hashed: %s -> %s\n", info.Name(), hash)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking the directory: %w", err)
	}

	return dlfiles, nil
}

// saveHashesToFile saves the file hashes to a JSON file.
func saveHashesToFile(dlfiles map[string]DLFile, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	if err := encoder.Encode(dlfiles); err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	fmt.Printf("File hashes saved to %s\n", outputPath)
	return nil

}
