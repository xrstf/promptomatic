// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package simple

import (
	"log"

	"go.xrstf.de/promptomatic/pkg/color"
	"go.xrstf.de/promptomatic/pkg/prompt"
	"go.xrstf.de/promptomatic/pkg/provider"
	"go.xrstf.de/promptomatic/pkg/sys"
	"go.xrstf.de/promptomatic/pkg/theme"
)

type simpleTheme struct{}

func New() theme.Theme {
	return &simpleTheme{}
}

func (t *simpleTheme) Render(e *sys.Environment) string {
	baseStyle := color.NewStyle().
		SetForeground(color.HiGreen).
		SetForegroundIf(color.HiRed, e.Elevated()).
		SetBold(true)

	b := prompt.NewBuilder(prompt.NewEmptyTransition())

	addSegment := func(builder func(*sys.Environment, *color.Style) (*prompt.Segment, error)) {
		segment, err := builder(e, baseStyle)
		if err != nil {
			log.Printf("Error: %v", err)
		} else if segment != nil {
			b.AddSegment(*segment)
		}
	}

	addSegment(t.whereami)
	addSegment(t.path)
	addSegment(t.git)
	addSegment(t.promptomatic)

	return b.Render()
}

func (t *simpleTheme) whereami(e *sys.Environment, style *color.Style) (*prompt.Segment, error) {
	return prompt.NewSegmentPtr("%m", *style), nil
}

func (t *simpleTheme) path(e *sys.Environment, _ *color.Style) (*prompt.Segment, error) {
	style := color.NewStyle().SetForeground(color.HiBlue).SetBold(true)

	return prompt.NewSegmentPtr("%~", *style), nil
}

func (t *simpleTheme) git(e *sys.Environment, _ *color.Style) (*prompt.Segment, error) {
	s, err := provider.NewGitStatus(e)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}

	style := color.NewStyle().
		SetForeground(color.Green).
		SetForegroundIf(color.Yellow, s.Dirty)

	text := "git:" + s.Branch

	if s.StagedChanges {
		text += "*"
	}

	return prompt.NewSegmentPtr(text, *style), nil
}

func (t *simpleTheme) promptomatic(e *sys.Environment, s *color.Style) (*prompt.Segment, error) {
	return prompt.NewSegmentPtr("%#", *s), nil
}
