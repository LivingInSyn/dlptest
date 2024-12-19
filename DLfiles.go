package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// hashFile calculates the SHA1 hash of a single file.
func hashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashing file: %w", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// hashDirectory calculates and prints the SHA1 hash of every file in a directory.
func hashDirectory(dirPath string) (map[string]DLFile, error) {
	var dlfiles = make(map[string]DLFile)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return nil // Continue walking even if there's an error with one file/directory
		}

		if !info.Mode().IsRegular() {
			return nil // Skip directories and other non-regular files
		}

		hash, err := hashFile(path)
		if err != nil {
			fmt.Printf("Error hashing %s: %v\n", path, err)
			return nil // Continue walking even if there's an error hashing one file
		}
		dlfile := DLFile{
			Name: info.Name(),
			Hash: hash,
		}
		dlfiles[info.Name()] = dlfile

		fmt.Printf("%s: %s\n", path, hash)
		return nil
	})

	return dlfiles, err
}
