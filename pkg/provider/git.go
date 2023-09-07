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

	if !isGitRepository(e.WorkingDirectory) {
		return nil, nil
	}

	branch, err := sys.RunCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		// swallow this error because it might just be not a valid git repo (yet)
		return nil, nil
	}

	changes, err := sys.RunCommand("git", "diff", "--name-only")
	if err != nil {
		return nil, err
	}

	staged, err := sys.RunCommand("git", "diff", "--name-only", "--cached")
	if err != nil {
		return nil, err
	}

	return &GitStatus{
		Branch:        branch,
		Dirty:         changes != "" || staged != "",
		StagedChanges: staged != "",
	}, nil
}

func isGitRepository(cwd string) bool {
	// we could for whatever reason not determine our current working dir
	if cwd == "" {
		return false
	}

	for cwd != "/" {
		gitDir := filepath.Join(cwd, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return true
		}

		cwd = filepath.Dir(cwd)
	}

	return false
}
