package cli

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// run executes the CLI in-process and returns combined stdout and stderr.
func run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var buf bytes.Buffer
	root := newRootCommand()
	root.SetOut(&buf)
	root.SetErr(&buf)
	root.SetArgs(args)
	err := root.ExecuteContext(context.Background())
	return buf.String(), err
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCommandsRegistered(t *testing.T) {
	want := map[string]bool{"version": false, "doctor": false, "analyse": false, "init": false, "skills": false, "github": false}
	for _, c := range newRootCommand().Commands() {
		if _, ok := want[c.Name()]; ok {
			if want[c.Name()] {
				t.Errorf("command %q registered twice", c.Name())
			}
			want[c.Name()] = true
		}
	}
	for name, ok := range want {
		if !ok {
			t.Errorf("command %q not registered", name)
		}
	}
}

func TestVersion(t *testing.T) {
	out, err := run(t, "version")
	if err != nil || !strings.HasPrefix(out, "Anthropica CLI ") {
		t.Fatalf("version: %q, %v", out, err)
	}
	if out2, _ := run(t, "--version"); out2 != out {
		t.Errorf("--version %q differs from version %q", out2, out)
	}
}

func TestUnknownCommand(t *testing.T) {
	if _, err := run(t, "bogus"); err == nil || !strings.Contains(err.Error(), "unknown command") {
		t.Fatalf("expected unknown command error, got %v", err)
	}
}

func TestDoctorHidesToken(t *testing.T) {
	t.Setenv("GITHUB_TOKEN", "ghp_secretvalue")
	out, err := run(t, "doctor")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "ghp_secretvalue") {
		t.Fatal("doctor printed the token")
	}
	if !strings.Contains(out, "GITHUB_TOKEN (hidden)") {
		t.Errorf("doctor did not report token source: %s", out)
	}
}

func TestAnalyse(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "pyproject.toml"), "")
	writeFile(t, filepath.Join(dir, "requirements.txt"), "")
	writeFile(t, filepath.Join(dir, "AGENTS.md"), "")
	writeFile(t, filepath.Join(dir, ".github", "workflows", "ci.yml"), "")

	out, err := run(t, "analyse", "--path", dir)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(out, "Python") != 1 {
		t.Errorf("Python should be reported once:\n%s", out)
	}
	for _, want := range []string{"GitHub Actions", "AGENTS.md"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in output:\n%s", want, out)
		}
	}
}

func TestAnalyseMissingPath(t *testing.T) {
	_, err := run(t, "analyse", "--path", filepath.Join(t.TempDir(), "nope"))
	if err == nil || !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("expected missing path error, got %v", err)
	}
}

func TestInitCreatesAndPreserves(t *testing.T) {
	dir := t.TempDir()
	existing := "# my own CLAUDE.md\n"
	writeFile(t, filepath.Join(dir, "CLAUDE.md"), existing)

	out, err := run(t, "init", "--path", dir, "--assistant", "all")
	if err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "CLAUDE.md"))
	if string(got) != existing {
		t.Fatalf("CLAUDE.md was overwritten: %q", got)
	}
	if !exists(filepath.Join(dir, "AGENTS.md")) || !exists(filepath.Join(dir, ".anthropica", "config.json")) {
		t.Fatalf("expected AGENTS.md and config.json to be created:\n%s", out)
	}
	if !strings.Contains(out, "2 created, 1 existing") {
		t.Errorf("unexpected summary:\n%s", out)
	}

	// A second run changes nothing.
	agents, _ := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	out, err = run(t, "init", "--path", dir, "--assistant", "all")
	if err != nil {
		t.Fatal(err)
	}
	again, _ := os.ReadFile(filepath.Join(dir, "AGENTS.md"))
	if !bytes.Equal(agents, again) || !strings.Contains(out, "0 created, 3 existing") {
		t.Errorf("second init was not a no-op:\n%s", out)
	}
}

func TestInitRejectsMissingPathAndBadAssistant(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "typo")
	if _, err := run(t, "init", "--path", missing); err == nil {
		t.Fatal("expected error for missing path")
	}
	if exists(missing) {
		t.Fatal("init created the missing project directory")
	}
	if _, err := run(t, "init", "--path", t.TempDir(), "--assistant", "nope"); err == nil {
		t.Fatal("expected error for invalid --assistant")
	}
}

func TestSkillsListFromSubdirectory(t *testing.T) {
	repo := t.TempDir()
	writeFile(t, filepath.Join(repo, "packs.json"), "{}")
	writeFile(t, filepath.Join(repo, "skills", "b-skill", "SKILL.md"), "---\nname: b-skill\ndescription: Second. Use when testing.\n---\n")
	writeFile(t, filepath.Join(repo, "skills", "a-skill", "SKILL.md"), "---\r\nname: a-skill\r\ndescription: First skill. Use when testing.\r\n---\r\n")
	writeFile(t, filepath.Join(repo, "plugins", "anthropica-qa", "skills", "a-skill", "SKILL.md"), "---\nname: a-skill\ndescription: First.\n---\n")
	sub := filepath.Join(repo, "cli", "internal")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatal(err)
	}

	out, err := run(t, "skills", "list", "--repo", sub)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "2 skills discovered") || strings.Index(out, "a-skill") > strings.Index(out, "b-skill") {
		t.Errorf("unexpected listing:\n%s", out)
	}
	if !strings.Contains(out, "First skill.") || strings.Contains(out, "Use when") {
		t.Errorf("expected first-sentence summaries:\n%s", out)
	}

	out, err = run(t, "skills", "list", "--repo", repo, "--pack", "anthropica-qa")
	if err != nil || !strings.Contains(out, "1 skills discovered") {
		t.Errorf("pack listing: %v\n%s", err, out)
	}
	if _, err := run(t, "skills", "list", "--repo", repo, "--pack", "../etc"); err == nil {
		t.Error("expected invalid pack name error")
	}
	if _, err := run(t, "skills", "list", "--repo", t.TempDir()); err == nil {
		t.Error("expected error outside an Anthropica repository")
	}
}

func TestParseRepo(t *testing.T) {
	valid := map[string]string{
		"jamal-gml/anthropica":                        "jamal-gml/anthropica",
		"https://github.com/jamal-gml/anthropica":     "jamal-gml/anthropica",
		"https://github.com/jamal-gml/anthropica.git": "jamal-gml/anthropica",
		"owner/repo.name_x":                           "owner/repo.name_x",
	}
	for in, want := range valid {
		o, r, err := parseRepo(in)
		if err != nil || o+"/"+r != want {
			t.Errorf("parseRepo(%q) = %s/%s, %v; want %s", in, o, r, err, want)
		}
	}
	for _, in := range []string{"", "owner", "owner/", "/repo", "a/b/c", "owner/..", "-owner/repo", "own er/repo", "owner/re?po", "owner/%2e%2e"} {
		if _, _, err := parseRepo(in); err == nil {
			t.Errorf("parseRepo(%q) should fail", in)
		}
	}
}

func TestGitHubInspect(t *testing.T) {
	var gotAuth []string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = append(gotAuth, r.Header.Get("Authorization"))
		switch r.URL.Path {
		case "/repos/octo/demo":
			w.Write([]byte(`{"full_name":"octo/demo","private":true,"default_branch":"main","html_url":"https://github.com/octo/demo"}`))
		case "/repos/octo/demo/git/trees/main":
			w.Write([]byte(`{"truncated":false,"tree":[{"path":"README.md","type":"blob"},{"path":"src","type":"tree"}]}`))
		case "/repos/octo/limited":
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()
	old := githubAPIBase
	githubAPIBase = srv.URL
	defer func() { githubAPIBase = old }()

	t.Setenv("GITHUB_TOKEN", "")
	t.Setenv("GH_TOKEN", "tok123")
	out, err := run(t, "github", "inspect", "octo/demo")
	if err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
	for _, want := range []string{"Visibility: private", "README.md", "src/"} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "tok123") {
		t.Error("token printed")
	}
	if len(gotAuth) == 0 || gotAuth[0] != "Bearer tok123" {
		t.Errorf("GH_TOKEN not sent: %v", gotAuth)
	}

	t.Setenv("GH_TOKEN", "")
	gotAuth = nil
	_, err = run(t, "github", "inspect", "octo/missing")
	if err == nil || !strings.Contains(err.Error(), "need GITHUB_TOKEN") {
		t.Errorf("expected 404 hint without token, got %v", err)
	}
	if len(gotAuth) == 0 || gotAuth[0] != "" {
		t.Errorf("Authorization header sent without a token: %v", gotAuth)
	}

	if _, err := run(t, "github", "inspect", "octo/limited"); err == nil || !strings.Contains(err.Error(), "rate limit") {
		t.Errorf("expected rate limit error, got %v", err)
	}
	if _, err := run(t, "github", "inspect", "not a repo"); err == nil {
		t.Error("expected invalid repository error")
	}
}
