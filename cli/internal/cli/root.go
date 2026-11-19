// Package cli implements the anthropica command-line interface.
package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// version is overridden at build time:
//
//	go build -ldflags "-X github.com/jamal-gml/anthropica/cli/internal/cli.version=1.2.0" .
var version = "dev"

// Version returns the build version, falling back to the module version
// recorded by `go install module@version`.
func Version() string {
	if version != "dev" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "anthropica",
		Short: "Anthropica: role-based AI assistant skills for your projects",
		Long: "Anthropica sets up projects for AI coding assistants and works with the\n" +
			"Anthropica skill library. It is provider-neutral: generated instructions\n" +
			"follow AGENTS.md, with optional assistant-specific files.",
		Version: Version(),
		// Errors are printed once by Execute; usage is only shown for usage errors.
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.SetVersionTemplate("Anthropica CLI {{.Version}}\n")

	root.AddCommand(
		newVersionCommand(),
		newDoctorCommand(),
		newAnalyseCommand(),
		newInitCommand(),
		newSkillsCommand(),
		newGitHubCommand(),
	)
	return root
}

// Execute runs the CLI and prints any error to stderr. Ctrl+C cancels
// in-flight work such as GitHub requests.
func Execute() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	root := newRootCommand()
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(root.ErrOrStderr(), "Error:", err)
		return err
	}
	return nil
}
