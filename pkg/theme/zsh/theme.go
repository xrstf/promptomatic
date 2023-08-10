// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package zsh

import (
	"go.xrstf.de/promptomatic/pkg/color"
	"go.xrstf.de/promptomatic/pkg/prompt"
	"go.xrstf.de/promptomatic/pkg/sys"
	"go.xrstf.de/promptomatic/pkg/theme"
)

type zshTheme struct{}

func New() theme.Theme {
	return &zshTheme{}
}

func (t *zshTheme) Render(e *sys.Environment) string {
	baseStyle := color.NewStyle().SetForegroundIf(color.HiRed, e.Elevated())

	b := prompt.NewBuilder(prompt.NewEmptyTransition())
	b.AddSegment(prompt.NewSegment("%m%#", *baseStyle))

	return b.Render()
}
