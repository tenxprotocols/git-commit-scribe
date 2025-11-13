package main

import (
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/tenxprotocols/git-commit-scribe/internal/cli"
)

var (
	version   = "dev"     // Set via ldflags during build
	commit    = "none"    // Set via ldflags during build
	buildDate = "unknown" // Set via ldflags during build
)

func main() {
	versionString := fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, buildDate)

	ctx := kong.Parse(&cli.CLI,
		kong.Name("gscribe"),
		kong.Description("AI-powered Git commit message generator"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": versionString,
		},
	)

	err := ctx.Run(&cli.Context{})
	ctx.FatalIfErrorf(err)
}
