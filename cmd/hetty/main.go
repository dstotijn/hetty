package main

import (
	"context"
	llog "log"
	"os"

	"go.uber.org/zap"

	"github.com/dstotijn/hetty/pkg/log"
)

func main() {
	hettyCmd, cfg := NewHettyCommand()

	if err := hettyCmd.Parse(os.Args[1:]); err != nil {
		llog.Fatalf("Failed to parse command line arguments: %v", err)
	}

	logger, err := log.NewZapLogger(cfg.verbose, cfg.jsonLogs)
	if err != nil {
		llog.Fatal(err)
	}
	defer logger.Sync()

	cfg.logger = logger

	if err := hettyCmd.Run(context.Background()); err != nil {
		logger.Fatal("Unexpected error running command.", zap.Error(err))
	}
}
