package web

import (
	"embed"
	"fmt"
	"io/fs"
)

//go:embed static
var EmbedFs embed.FS

// PrintEmbeddedFiles prints out all files in the embedded filesystem
func PrintEmbeddedFiles() {
	fmt.Println("Contents of embedded filesystem:")
	fs.WalkDir(EmbedFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fmt.Printf("- %s\n", path)
		return nil
	})
}