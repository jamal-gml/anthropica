package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

const universalInstructions = `# AGENTS.md

## Anthropica Development Principles

- Establish outcomes and acceptance criteria before implementation.
- Inspect project instructions, Git status and relevant files.
- Use targeted searches rather than dumping entire repositories.
- Separate verified facts from assumptions.
- Preserve existing architecture and conventions.
- Never claim tests passed unless actually executed.
- Never expose credentials or weaken security controls.
- Do not perform destructive actions without approval.
- Use only the relevant skills for each task.
- Finish substantial work with a summary of changes,
  validation, risks and next steps.
`

const claudeInstructions = `# CLAUDE.md

Follow the repository-wide instructions in AGENTS.md.

Apply the relevant Anthropica skill for the current task.

Inspect existing architecture before making changes.

Preserve existing project conventions.

Never claim validation that was not performed.
`

var assistantModes = map[string]bool{"universal": true, "claude": true, "all": true}

func newInitCommand() *cobra.Command {
	var path string
	var assistant string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialise Anthropica in a project (never overwrites existing files)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !assistantModes[assistant] {
				return fmt.Errorf("--assistant must be universal, claude or all (got %q)", assistant)
			}
			// Do not create the project directory: a typo in --path should fail, not scaffold elsewhere.
			root, err := existingDir(path)
			if err != nil {
				return err
			}
			out := cmd.OutOrStdout()

			configDir := filepath.Join(root, ".anthropica")
			if err := os.MkdirAll(configDir, 0o755); err != nil {
				return fmt.Errorf("cannot create %s: %w", configDir, err)
			}

			config := map[string]any{
				"schemaVersion": 1,
				"tool":          "anthropica",
				"assistantMode": assistant,
				"createdAt":     time.Now().UTC().Format(time.RFC3339),
			}
			data, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				return err
			}

			files := []struct {
				path    string
				content []byte
				want    bool
			}{
				{filepath.Join(configDir, "config.json"), append(data, '\n'), true},
				{filepath.Join(root, "AGENTS.md"), []byte(universalInstructions), assistant != "claude"},
				{filepath.Join(root, "CLAUDE.md"), []byte(claudeInstructions), assistant != "universal"},
			}

			created, skipped := 0, 0
			for _, f := range files {
				if !f.want {
					continue
				}
				wrote, err := writeIfMissing(out, f.path, f.content)
				if err != nil {
					return err
				}
				if wrote {
					created++
				} else {
					skipped++
				}
			}

			fmt.Fprintf(out, "\nAnthropica initialised: %d created, %d existing file(s) left unchanged.\n", created, skipped)
			return nil
		},
	}

	cmd.Flags().StringVarP(&path, "path", "p", ".", "Project directory (must already exist)")
	cmd.Flags().StringVar(&assistant, "assistant", "universal", "universal (AGENTS.md), claude (CLAUDE.md) or all")
	return cmd
}

// writeIfMissing creates path with content, or leaves an existing file
// untouched. It reports whether the file was written.
func writeIfMissing(out io.Writer, path string, content []byte) (bool, error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if errors.Is(err, fs.ErrExist) {
		fmt.Fprintln(out, "[SKIP]", path, "already exists")
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("cannot create %s: %w", path, err)
	}

	_, writeErr := file.Write(content)
	closeErr := file.Close()
	if err := errors.Join(writeErr, closeErr); err != nil {
		// Remove the partial file we created so a re-run can succeed.
		_ = os.Remove(path)
		return false, fmt.Errorf("cannot write %s: %w", path, err)
	}

	fmt.Fprintln(out, "[CREATED]", path)
	return true, nil
}
