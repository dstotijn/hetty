package main

import (
	"flag"

	"go.uber.org/zap"
)

// Config represents the global configuration shared amongst all commands.
type Config struct {
	verbose  bool
	jsonLogs bool
	logger   *zap.Logger
}

// RegisterFlags registers the flag fields into the provided flag.FlagSet. This
// helper function allows subcommands to register the root flags into their
// flagsets, creating "global" flags that can be passed after any subcommand at
// the commandline.
func (cfg *Config) RegisterFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cfg.verbose, "verbose", false, "Enable verbose logging.")
	fs.BoolVar(&cfg.jsonLogs, "json", false, "Encode logs as JSON, instead of pretty/human readable output.")
}
