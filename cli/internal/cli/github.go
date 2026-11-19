package cli

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// githubAPIBase is a variable so tests can point it at a local server.
var githubAPIBase = "https://api.github.com"

const maxGitHubResponse = 16 << 20 // 16 MiB

var (
	ownerPattern = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9-]{0,38})$`)
	repoPattern  = regexp.MustCompile(`^[A-Za-z0-9._-]{1,100}$`)
)

type githubRepository struct {
	FullName      string `json:"full_name"`
	Private       bool   `json:"private"`
	DefaultBranch string `json:"default_branch"`
	Description   string `json:"description"`
	HTMLURL       string `json:"html_url"`
}

type githubTree struct {
	Truncated bool `json:"truncated"`
	Tree      []struct {
		Path string `json:"path"`
		Type string `json:"type"`
	} `json:"tree"`
}

func newGitHubCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "github",
		Short: "GitHub repository tools (read-only)",
	}

	var token string
	inspect := &cobra.Command{
		Use:   "inspect OWNER/REPOSITORY",
		Short: "Inspect GitHub repository metadata",
		Long: "Read-only inspection of a GitHub repository's metadata and root entries.\n" +
			"Uses GITHUB_TOKEN or GH_TOKEN when set; public repositories work without one.",
		Example: "  anthropica github inspect jamal-gml/anthropica\n" +
			"  anthropica github inspect https://github.com/jamal-gml/anthropica",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, repo, err := parseRepo(args[0])
			if err != nil {
				return err
			}
			out, errOut := cmd.OutOrStdout(), cmd.ErrOrStderr()

			if token != "" {
				fmt.Fprintln(errOut, "[WARNING] --token can leak into shell history; prefer the GITHUB_TOKEN environment variable.")
			} else if name := githubTokenSource(); name != "" {
				token = os.Getenv(name)
			} else {
				fmt.Fprintln(errOut, "[INFO] No GitHub token set; only public repositories are visible and rate limits are low.")
			}

			ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
			defer cancel()

			var meta githubRepository
			repoURL := fmt.Sprintf("%s/repos/%s/%s", githubAPIBase, url.PathEscape(owner), url.PathEscape(repo))
			if err := githubGet(ctx, token, repoURL, &meta); err != nil {
				return err
			}

			visibility := "public"
			if meta.Private {
				visibility = "private"
			}
			fmt.Fprintf(out, "\nRepository: %s\n", meta.FullName)
			fmt.Fprintf(out, "Visibility: %s\n", visibility)
			fmt.Fprintf(out, "Default branch: %s\n", meta.DefaultBranch)
			fmt.Fprintf(out, "URL: %s\n", meta.HTMLURL)
			if meta.Description != "" {
				fmt.Fprintf(out, "Description: %s\n", meta.Description)
			}

			if meta.DefaultBranch == "" {
				fmt.Fprintln(errOut, "\n[WARNING] Repository has no default branch (empty repository?); skipping tree.")
				return nil
			}

			// Without recursive=1 the API returns only root entries.
			var tree githubTree
			treeURL := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s",
				githubAPIBase, url.PathEscape(owner), url.PathEscape(repo), url.PathEscape(meta.DefaultBranch))
			if err := githubGet(ctx, token, treeURL, &tree); err != nil {
				return fmt.Errorf("repository metadata loaded, but listing root entries failed: %w", err)
			}

			fmt.Fprintln(out, "\nRoot entries:")
			for _, item := range tree.Tree {
				suffix := ""
				if item.Type == "tree" {
					suffix = "/"
				}
				fmt.Fprintf(out, " - %s%s\n", item.Path, suffix)
			}
			if tree.Truncated {
				fmt.Fprintln(out, " (list truncated by GitHub)")
			}

			fmt.Fprintln(out, "\nRead-only inspection completed.")
			return nil
		},
	}
	inspect.Flags().StringVar(&token, "token", "", "GitHub token (prefer the GITHUB_TOKEN environment variable)")

	cmd.AddCommand(inspect)
	return cmd
}

// parseRepo accepts OWNER/REPO, optionally as a github.com URL or with a .git
// suffix, and validates both parts against GitHub's naming rules.
func parseRepo(input string) (owner, repo string, err error) {
	s := strings.TrimSpace(input)
	for _, prefix := range []string{"https://github.com/", "http://github.com/", "github.com/"} {
		if strings.HasPrefix(strings.ToLower(s), prefix) {
			s = s[len(prefix):]
			break
		}
	}
	s = strings.TrimSuffix(strings.TrimSuffix(s, "/"), ".git")

	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository %q: use OWNER/REPOSITORY", input)
	}
	owner, repo = parts[0], parts[1]
	if !ownerPattern.MatchString(owner) {
		return "", "", fmt.Errorf("invalid repository owner %q", owner)
	}
	if !repoPattern.MatchString(repo) || repo == "." || repo == ".." {
		return "", "", fmt.Errorf("invalid repository name %q", repo)
	}
	return owner, repo, nil
}

func githubGet(ctx context.Context, token, apiURL string, target any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "anthropica-cli/"+Version())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("GitHub request failed: %w", err)
	}
	defer resp.Body.Close()

	body := io.LimitReader(resp.Body, maxGitHubResponse)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, body)
		return githubStatusError(resp, token != "")
	}
	if err := json.NewDecoder(body).Decode(target); err != nil {
		return fmt.Errorf("could not decode GitHub response: %w", err)
	}
	return nil
}

func githubStatusError(resp *http.Response, authenticated bool) error {
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("GitHub returned %s: the token is invalid or expired", resp.Status)
	case http.StatusForbidden, http.StatusTooManyRequests:
		if resp.Header.Get("X-RateLimit-Remaining") == "0" {
			hint := "try again later"
			if !authenticated {
				hint = "set GITHUB_TOKEN for a higher limit"
			}
			return fmt.Errorf("GitHub API rate limit exceeded; %s", hint)
		}
		return fmt.Errorf("GitHub returned %s: the token lacks permission for this repository", resp.Status)
	case http.StatusNotFound:
		if authenticated {
			return fmt.Errorf("GitHub returned %s: repository not found, or the token cannot access it", resp.Status)
		}
		return fmt.Errorf("GitHub returned %s: repository not found (private repositories need GITHUB_TOKEN)", resp.Status)
	default:
		return fmt.Errorf("GitHub returned %s", resp.Status)
	}
}
