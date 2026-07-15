package main

import (
	"fmt"
	"log"

	"github.com/RyanConnell/dotfiles/dotter/internal/commands"
	"github.com/alecthomas/kong"
)

type Flags struct {
	Config  string `help:"Path to config file" default:"config.yaml" short:"c" type:"path"`
	Apps    string `help:"Path to apps folder" default:"apps" short:"a" type:"path"`
	Install struct {
		Output string `help:"Path to output folder" short:"o" type:"path"`
	} `cmd:"install" help:"Install templates"`
}

func main() {
	var flags Flags
	ctx := kong.Parse(&flags,
		kong.Name("dotter"),
		kong.Description("Dotfile templater, installer, and monitor"),
		kong.UsageOnError(),
	)

	var err error
	switch ctx.Command() {
	case "install":
		cmd := commands.NewInstaller(flags.Config, flags.Apps, flags.Install.Output)
		err = cmd.Run()
	default:
		err = fmt.Errorf("unknown subcommand")
	}
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
