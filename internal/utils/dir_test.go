package utils

import (
	"os"
	"testing"
)

// Testing Dir function
func TestDir(t *testing.T) {
	dir := Dir()
	if dir == "" {
		t.Error("Expected a directory path, got an empty string")
	}
}

// Testing UploadDir function
func TestUploadDir(t *testing.T) {
	_, fullPath := UploadDir()
	_, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		t.Errorf("Expected directory to exist at path: %s", fullPath)
	}

}
