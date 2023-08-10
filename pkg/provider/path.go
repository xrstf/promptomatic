// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package provider

import (
	"strings"

	"go.xrstf.de/promptomatic/pkg/sys"
)

const (
	gopher = "\ue626"
)

type PathStatus struct {
	Path   string
	Symbol string
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

	result = strings.TrimSuffix(result, "/")

	return &PathStatus{
		Path:   result,
		Symbol: symbol,
	}, nil
}
