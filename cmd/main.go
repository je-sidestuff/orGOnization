package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/je-sidestuff/orgonization/templates"
)

func handleBuiltinTemplate(builtin string) {
	switch builtin {
	case "weekly":
		handleWeeklyTemplate()
	default:
		fmt.Printf("Unknown builtin template: %s\n", builtin)
		os.Exit(1)
	}
}

func handleWeeklyTemplate() {
	err := templates.PrintWeekdays(time.Now())
	if err != nil {
		fmt.Printf("Error generating weekly template: %v\n", err)
		os.Exit(1)
	}
}

func main() {

	// Define the main command and subcommands as separate flag sets
	var mainCmd = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Usage comment explaining how to use the CLI tool
	usage := `

	The goorganizethings module is a Go executable capable of running one-off CLI commands, as a local-only agent, or as a server.

	Usage:
	%s <subcommand> [flags]

	Available Subcommands:
	template: Outputs a specified populated template or performs a one-off processing of an IO file tree
	agent:    Runs continuously as an agent, watching an IO file tree and performing updates according to its configuration
	server:   Runs as an agent backed by a server, allowing a client to make updates while IO file tree processing is conducted

	Flags:
	-h, --help    show help message

	**Subcommand-specific flags are available. Run the subcommand with the "-h" flag for details.**

`

	// Help flag
	var help bool
	mainCmd.BoolVar(&help, "h", false, "show help message")
	mainCmd.BoolVar(&help, "help", false, "show help message (shorthand)")

	// Print usage message and exit if help flag is set or arguments are incorrect
	if help || len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}

	// Define template subcommand and its flags
	templateCmd := flag.NewFlagSet("template", flag.ExitOnError)
	var buildTarget string
	templateCmd.StringVar(&buildTarget, "file", "", "Top level file to run templating on.")
	var buildVerbose bool
	templateCmd.BoolVar(&buildVerbose, "verbose", false, "Enable verbose output")
	var templateFlag1 string
	templateCmd.StringVar(&templateFlag1, "recurse", "", "Whether to recurse templating to generated files.")
	var builtin string
	templateCmd.StringVar(&builtin, "builtin", "", "Builtin template to use (optional)")

	// Parse command line arguments with main command
	if err := mainCmd.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	subcmd := os.Args[1]

	// Parse subcommand specific flags based on chosen subcommand
	switch subcmd {
	case "template":
		if err := templateCmd.Parse(os.Args[2:]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Handle builtin argument (now optional)
		if builtin != "" {
			handleBuiltinTemplate(builtin)
		}

	default:
		fmt.Fprintf(os.Stderr, "Invalid subcommand: %s", os.Args[1])
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		os.Exit(1)
	}
}
