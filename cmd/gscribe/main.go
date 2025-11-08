package main

import (
	"github.com/alecthomas/kong"
	"github.com/tenxprotocols/git-commit-scribe/internal/cli"
)

var version = "dev"

func main() {
	ctx := kong.Parse(&cli.CLI,
		kong.Name("gscribe"),
		kong.Description("AI-powered Git commit message generator"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version": version,
		},
	)

	err := ctx.Run(&cli.Context{})
	ctx.FatalIfErrorf(err)
}
