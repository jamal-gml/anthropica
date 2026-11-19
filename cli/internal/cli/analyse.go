package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

// manifests maps a root-level file to the technology it indicates.
var manifests = []struct{ file, technology string }{
	{"go.mod", "Go"},
	{"package.json", "Node.js / JavaScript / TypeScript"},
	{"pyproject.toml", "Python"},
	{"requirements.txt", "Python"},
	{"Cargo.toml", "Rust"},
	{"pom.xml", "Java / Maven"},
	{"build.gradle", "Java / Gradle"},
	{"build.gradle.kts", "Java / Gradle"},
	{"Gemfile", "Ruby"},
	{"composer.json", "PHP"},
	{"Dockerfile", "Docker"},
	{".github/workflows", "GitHub Actions"},
}

// instructionFiles are AI assistant instruction files, provider-neutral first.
var instructionFiles = []string{
	"AGENTS.md",
	"CLAUDE.md",
	"GEMINI.md",
	".github/copilot-instructions.md",
	"COPILOT.md",
	".anthropica/config.json",
}

func newAnalyseCommand() *cobra.Command {
	var path string

	cmd := &cobra.Command{
		Use:     "analyse",
		Aliases: []string{"analyze"},
		Short:   "Analyse a local project",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			root, err := existingDir(path)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()

			entries, err := os.ReadDir(root)
			if err != nil {
				return fmt.Errorf("cannot read %s: %w", root, err)
			}

			fmt.Fprintf(out, "\nProject: %s\n\nRoot files:\n", root)
			names := make([]string, 0, len(entries))
			for _, entry := range entries {
				name := entry.Name()
				if entry.IsDir() {
					name += string(filepath.Separator)
				}
				names = append(names, name)
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Fprintln(out, " -", name)
			}

			fmt.Fprintln(out, "\nDetected technologies:")
			seen := map[string]bool{}
			for _, m := range manifests {
				if seen[m.technology] || !exists(filepath.Join(root, filepath.FromSlash(m.file))) {
					continue
				}
				seen[m.technology] = true
				fmt.Fprintf(out, " - %s (%s)\n", m.technology, m.file)
			}
			if len(seen) == 0 {
				fmt.Fprintln(out, " - No recognised manifest found")
			}

			fmt.Fprintln(out, "\nAI instruction files:")
			found := false
			for _, file := range instructionFiles {
				if exists(filepath.Join(root, filepath.FromSlash(file))) {
					fmt.Fprintln(out, " -", file)
					found = true
				}
			}
			if !found {
				fmt.Fprintln(out, " - None found (run `anthropica init` to create AGENTS.md)")
			}

			fmt.Fprintln(out, "\nAnalysis completed.")
			fmt.Fprintln(out, "Note: this is a metadata scan, not a full code audit.")
			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "Project directory")
	return cmd
}

// existingDir resolves path to an absolute directory and errors clearly if it
// does not exist or is not a directory.
func existingDir(path string) (string, error) {
	root, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("invalid path %q: %w", path, err)
	}
	info, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("path does not exist: %s", root)
		}
		return "", fmt.Errorf("cannot access %s: %w", root, err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("path is not a directory: %s", root)
	}
	return root, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
