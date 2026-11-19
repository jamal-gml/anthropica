package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

func newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check system prerequisites",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			fmt.Fprintf(out, "Anthropica CLI: %s\n", Version())
			fmt.Fprintf(out, "Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)

			for _, tool := range []struct{ bin, name string }{
				{"git", "Git"},
				{"gh", "GitHub CLI"},
			} {
				if path, err := exec.LookPath(tool.bin); err == nil {
					fmt.Fprintf(out, "[OK] %s installed (%s)\n", tool.name, path)
				} else {
					fmt.Fprintf(out, "[INFO] %s not found on PATH\n", tool.name)
				}
			}

			// Never print the token itself.
			if name := githubTokenSource(); name != "" {
				fmt.Fprintf(out, "[OK] GitHub token configured via %s (hidden)\n", name)
			} else {
				fmt.Fprintln(out, "[INFO] GitHub token not configured (set GITHUB_TOKEN or GH_TOKEN for private repositories and higher rate limits)")
			}
		},
	}
}

// githubTokenSource returns the name of the environment variable holding a
// GitHub token, or "" if none is set.
func githubTokenSource() string {
	for _, name := range []string{"GITHUB_TOKEN", "GH_TOKEN"} {
		if os.Getenv(name) != "" {
			return name
		}
	}
	return ""
}
