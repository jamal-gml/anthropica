package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type skill struct {
	Name        string
	Description string
}

func newSkillsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Work with Anthropica skills",
	}

	var repo, pack string
	list := &cobra.Command{
		Use:   "list",
		Short: "List available skills",
		Long: "List skills from an Anthropica repository checkout. The repository root is\n" +
			"found by searching upwards from --repo. Use --pack to list one role pack.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			start, err := existingDir(repo)
			if err != nil {
				return err
			}
			root, err := findRepoRoot(start)
			if err != nil {
				return err
			}
			dir, err := skillsDir(root, pack)
			if err != nil {
				return err
			}
			skills, err := loadSkills(dir)
			if err != nil {
				return err
			}

			out := cmd.OutOrStdout()
			label := "all skills"
			if pack != "" {
				label = "pack " + pack
			}
			fmt.Fprintf(out, "\nANTHROPICA SKILLS (%s)\n\n", label)
			tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
			for _, s := range skills {
				fmt.Fprintf(tw, " %s\t%s\n", s.Name, summary(s.Description))
			}
			if err := tw.Flush(); err != nil {
				return err
			}
			fmt.Fprintf(out, "\n%d skills discovered in %s\n", len(skills), dir)
			return nil
		},
	}
	list.Flags().StringVarP(&repo, "repo", "r", ".", "Anthropica repository root or any directory inside it")
	list.Flags().StringVar(&pack, "pack", "", "Role pack to list, e.g. anthropica-devops (default: all skills)")

	cmd.AddCommand(list)
	return cmd
}

// findRepoRoot walks up from start to the directory that contains the
// Anthropica skill library.
func findRepoRoot(start string) (string, error) {
	dir := start
	for {
		if exists(filepath.Join(dir, "packs.json")) ||
			exists(filepath.Join(dir, ".claude-plugin", "marketplace.json")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no Anthropica repository found at or above %s (looked for packs.json or .claude-plugin/marketplace.json)", start)
		}
		dir = parent
	}
}

// skillsDir returns the folder holding SKILL.md files: a role pack's plugin,
// the single-source skills/ folder, or the legacy single-plugin layout.
func skillsDir(root, pack string) (string, error) {
	if pack != "" {
		if strings.ContainsAny(pack, `/\.`) {
			return "", fmt.Errorf("invalid pack name %q", pack)
		}
		dir := filepath.Join(root, "plugins", pack, "skills")
		if !exists(dir) {
			return "", fmt.Errorf("pack %q not found in %s", pack, filepath.Join(root, "plugins"))
		}
		return dir, nil
	}
	for _, dir := range []string{
		filepath.Join(root, "skills"),
		filepath.Join(root, "plugins", "anthropica", "skills"),
	} {
		if exists(dir) {
			return dir, nil
		}
	}
	return "", fmt.Errorf("no skills folder found in %s", root)
}

func loadSkills(dir string) ([]skill, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot read skills in %s: %w", dir, err)
	}
	var skills []skill
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		file := filepath.Join(dir, entry.Name(), "SKILL.md")
		if !exists(file) {
			continue
		}
		meta, err := readFrontmatter(file)
		if err != nil {
			return nil, err
		}
		name := meta["name"]
		if name == "" {
			name = entry.Name()
		}
		skills = append(skills, skill{Name: name, Description: meta["description"]})
	}
	sort.Slice(skills, func(i, j int) bool { return skills[i].Name < skills[j].Name })
	return skills, nil
}

// readFrontmatter reads simple top-level "key: value" pairs from a YAML
// frontmatter block. Nested YAML is ignored; it is not needed for listing.
func readFrontmatter(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer f.Close()

	meta := map[string]string{}
	sc := bufio.NewScanner(f)
	for i := 0; sc.Scan(); i++ {
		line := strings.TrimRight(sc.Text(), "\r")
		if i == 0 {
			if line != "---" {
				return meta, nil
			}
			continue
		}
		if line == "---" {
			break
		}
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			continue
		}
		if key, value, ok := strings.Cut(line, ":"); ok {
			meta[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("cannot read %s: %w", path, err)
	}
	return meta, nil
}

// summary returns the first sentence of a description, capped for a table.
func summary(desc string) string {
	if i := strings.Index(desc, ". "); i > 0 {
		desc = desc[:i+1]
	}
	const max = 100
	if r := []rune(desc); len(r) > max {
		desc = string(r[:max-1]) + "…"
	}
	return desc
}
