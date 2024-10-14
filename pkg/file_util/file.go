package file_util

import (
	"fmt"
	"os"
	"path/filepath"
)

// Create creates a file in the specified directory.
// If the directory does not exist, it creates the directory first.
// It then writes the provided content to the created file.
func Create(directory, filename, content string) error {
	// Check if the directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.MkdirAll(directory, 0755) // Use 0755 permissions for directories
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Create the full path for the file
	filePath := filepath.Join(directory, filename)

	// Create a new file
	file, err := os.Create(filePath) // os.Create creates or truncates the file
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close() // Close the file when we're done

	// Write content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
