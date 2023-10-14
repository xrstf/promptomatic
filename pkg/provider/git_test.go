// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestGetGitWorkingCopyRoot(t *testing.T) {
	testcases := []struct {
		dirs     []string
		pwd      string
		expected string
	}{
		{
			dirs:     []string{"foo", "foo/.git/hooks", "foo/bar/bla"},
			pwd:      "foo",
			expected: "foo",
		},
		{
			dirs:     []string{"foo", "foo/.git/hooks", "foo/bar/bla"},
			pwd:      "foo/bar",
			expected: "foo",
		},
		{
			dirs:     []string{"foo", "foo/.git/hooks", "foo/bar/bla"},
			pwd:      "foo/.git",
			expected: "foo",
		},
		{
			dirs:     []string{"foo", "foo/.git/hooks", "foo/bar/bla"},
			pwd:      "foo/.git/hooks",
			expected: "foo",
		},
		{
			dirs: []string{
				"foo/bla/.git/hooks",
				"foo/blub/.git/hooks",
			},
			pwd:      "foo",
			expected: "",
		},
		{
			dirs: []string{
				"foo/bla/.git/hooks",
				"foo/blub/.git/hooks",
			},
			pwd:      "foo/bla/.git/hooks",
			expected: "foo/bla",
		},
		{
			dirs: []string{
				"foo/bla/.git/hooks",
				"foo/blub/.git/hooks",
			},
			pwd:      "foo/blub/.git/hooks",
			expected: "foo/blub",
		},
		{
			dirs: []string{
				"foo/bla/.git/hooks",
				"foo/bla/thing/subrepo/.git",
			},
			pwd:      "foo/bla/thing",
			expected: "foo/bla",
		},
		{
			dirs: []string{
				"foo/bla/.git/hooks",
				"foo/bla/thing/subrepo/.git",
			},
			pwd:      "foo/bla/thing/subrepo",
			expected: "foo/bla/thing/subrepo",
		},
		{
			dirs: []string{
				"foo/bla/.git/hooks",
				"foo/bla/thing/subrepo/.git",
			},
			pwd:      "foo/bla/thing/subrepo/.git",
			expected: "foo/bla/thing/subrepo",
		},
	}

	for i, testcase := range testcases {
		t.Run(fmt.Sprintf("testcase %d", i), func(t *testing.T) {
			rootDir, err := setupTestcase(testcase.dirs)
			defer func() {
				if rootDir != "" {
					os.RemoveAll(rootDir)
				}
			}()

			if err != nil {
				t.Fatalf("Failed to setup: %v", err)
			}

			repoRoot := getGitWorkingCopyRoot(filepath.Join(rootDir, testcase.pwd))
			if repoRoot != "" {
				repoRoot, _ = filepath.Rel(rootDir, repoRoot)
			}

			if repoRoot != testcase.expected {
				t.Fatalf("Expected git root to be %q, but got %q.", testcase.expected, repoRoot)
			}
		})
	}
}

func setupTestcase(dirs []string) (string, error) {
	rootDir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}

	for _, dir := range dirs {
		fullDir := filepath.Join(rootDir, dir)

		if err := os.MkdirAll(fullDir, 0755); err != nil {
			return "", err
		}
	}

	return rootDir, nil
}
