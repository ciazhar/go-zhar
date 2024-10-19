package web

import (
	"context"
	"embed"
	"fmt"
	"io/fs"

	"github.com/ciazhar/go-start-small/pkg/logger"
)

//go:embed static
var EmbedFs embed.FS

// PrintEmbeddedFiles prints out all files in the embedded filesystem
func PrintEmbeddedFiles() {
	logger.LogInfo(context.Background(), "Contents of embedded filesystem", nil)
	fs.WalkDir(EmbedFs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		logger.LogInfo(context.Background(), fmt.Sprintf("- %s\n", path), nil)
		return nil
	})
}