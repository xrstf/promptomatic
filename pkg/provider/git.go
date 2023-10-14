// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package provider

import (
	"os"
	"path/filepath"

	"go.xrstf.de/promptomatic/pkg/sys"
)

type GitStatus struct {
	Branch        string
	Dirty         bool
	StagedChanges bool
}

func NewGitStatus(e *sys.Environment) (*GitStatus, error) {
	if !sys.CommandExists("git") {
		return nil, nil
	}

	repoRoot := getGitWorkingCopyRoot(e.WorkingDirectory)
	if repoRoot == "" {
		return nil, nil
	}

	branch, err := sys.RunCommand(repoRoot, "git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		// swallow this error because it might just be not a valid git repo (yet)
		return nil, nil
	}

	changes, err := sys.RunCommand(repoRoot, "git", "diff", "--name-only")
	if err != nil {
		return nil, err
	}

	staged, err := sys.RunCommand(repoRoot, "git", "diff", "--name-only", "--cached")
	if err != nil {
		return nil, err
	}

	return &GitStatus{
		Branch:        branch,
		Dirty:         changes != "" || staged != "",
		StagedChanges: staged != "",
	}, nil
}

func getGitWorkingCopyRoot(cwd string) string {
	// we could for whatever reason not determine our current working dir
	if cwd == "" {
		return ""
	}

	for cwd != "/" {
		// we're not just inside a Git working copy, we're inside its .git directory
		if filepath.Base(cwd) == ".git" {
			return filepath.Dir(cwd)
		}

		gitDir := filepath.Join(cwd, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return cwd
		}

		cwd = filepath.Dir(cwd)
	}

	return ""
}
