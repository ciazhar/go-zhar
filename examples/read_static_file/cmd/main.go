package main

import (
	"context"
	"github.com/ciazhar/go-start-small/examples/read_static_file/web"
	"github.com/ciazhar/go-start-small/pkg/logger"
)

func main() {
	file, err := web.EmbedFs.ReadFile("static/index.html")
	if err != nil {
		logger.LogAndReturnError(context.Background(), err, "Failed to read file", nil)
		return
	}

	logger.LogInfo(context.Background(), string(file), nil)
}
