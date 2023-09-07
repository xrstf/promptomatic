// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package provider

import (
	"os"
	"strings"

	"go.xrstf.de/promptomatic/pkg/sys"
)

const (
	gopher = "\ue626"
)

type PathStatus struct {
	Path   string
	Symbol string
	Exists bool
}

func NewPathStatus(e *sys.Environment) (*PathStatus, error) {
	path := e.WorkingDirectory
	result := path
	symbol := ""

	if strings.HasPrefix(path, e.HomeDirectory) {
		result = strings.Replace(path, e.HomeDirectory, "~", 1)
	}

	if e.GoPath != "" && strings.HasPrefix(path, e.GoPath) {
		symbol = gopher
		result = strings.Replace(path, e.GoPath, "", 1)
	}

	exists := true
	if _, err := os.Stat(e.WorkingDirectory); err != nil {
		exists = false
	}

	result = strings.TrimSuffix(result, "/")

	return &PathStatus{
		Path:   result,
		Symbol: symbol,
		Exists: exists,
	}, nil
}
