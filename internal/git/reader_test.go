package git

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommitDiffRoundTripRebuildsExactCommits(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()

	runTestGit(t, src, "init")
	runTestGit(t, src, "checkout", "-b", "main")
	runTestGit(t, src, "config", "user.name", "Alice")
	runTestGit(t, src, "config", "user.email", "alice@example.test")

	if err := os.WriteFile(filepath.Join(src, "README.md"), []byte("hello\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, src, "add", "README.md")
	runTestGitEnv(t, src, fixedCommitEnv(1), "commit", "-m", "initial commit")

	if err := os.WriteFile(filepath.Join(src, "README.md"), []byte("hello\nworld\n"), 0644); err != nil {
		t.Fatal(err)
	}
	runTestGit(t, src, "add", "README.md")
	runTestGitEnv(t, src, fixedCommitEnv(2), "commit", "-m", "second commit")

	runTestGit(t, dst, "init")
	if err := os.MkdirAll(filepath.Join(dst, ".bit"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dst, ".bit", "config.json"), []byte("{}\n"), 0600); err != nil {
		t.Fatal(err)
	}

	commits := strings.Fields(string(runTestGit(t, src, "rev-list", "--reverse", "HEAD")))
	if len(commits) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(commits))
	}

	for _, commit := range commits {
		info, err := ReadCommit(src, commit)
		if err != nil {
			t.Fatal(err)
		}
		diff, err := ExtractCommitDiff(src, info)
		if err != nil {
			t.Fatal(err)
		}
		m := ManifestForCommit("main", info, "bafyDiff")
		if err := ApplyCommitDiff(dst, m, diff); err != nil {
			t.Fatal(err)
		}
		got := strings.TrimSpace(string(runTestGit(t, dst, "rev-parse", "HEAD")))
		if got != commit {
			t.Fatalf("rebuilt commit mismatch: got %s, want %s", got, commit)
		}
	}
}

func fixedCommitEnv(n int) []string {
	stamp := "@" + []string{"1700000000", "1700000100"}[n-1] + " +0000"
	return []string{
		"GIT_AUTHOR_NAME=Alice",
		"GIT_AUTHOR_EMAIL=alice@example.test",
		"GIT_AUTHOR_DATE=" + stamp,
		"GIT_COMMITTER_NAME=Alice",
		"GIT_COMMITTER_EMAIL=alice@example.test",
		"GIT_COMMITTER_DATE=" + stamp,
	}
}

func runTestGit(t *testing.T, dir string, args ...string) []byte {
	t.Helper()
	return runTestGitEnv(t, dir, nil, args...)
}

func runTestGitEnv(t *testing.T, dir string, env []string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("git %s failed: %v: %s", strings.Join(args, " "), err, stderr.String())
	}
	return out.Bytes()
}
