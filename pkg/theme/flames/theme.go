// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package flames

import (
	"log"

	"go.xrstf.de/promptomatic/pkg/color"
	"go.xrstf.de/promptomatic/pkg/prompt"
	"go.xrstf.de/promptomatic/pkg/provider"
	"go.xrstf.de/promptomatic/pkg/sys"
	"go.xrstf.de/promptomatic/pkg/theme"
)

type flamesTheme struct{}

func New() theme.Theme {
	return &flamesTheme{}
}

func (t *flamesTheme) Render(e *sys.Environment) string {
	b := prompt.NewBuilder(prompt.NewSymbolTransition(" "))

	addSegment := func(builder func(*sys.Environment) (*prompt.Segment, error)) {
		segment, err := builder(e)
		if err != nil {
			log.Printf("Error: %v", err)
		} else if segment != nil {
			b.AddSegment(*segment)
		}
	}

	addSegment(t.whereami)
	addSegment(t.path)
	addSegment(t.git)
	addSegment(t.kubernetes)

	return b.Render()
}

func (t *flamesTheme) whereami(e *sys.Environment) (*prompt.Segment, error) {
	s, err := provider.NewWhoamiStatus(e)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}

	style := color.NewStyle().
		SetBackground(color.Green).
		SetBackgroundIf(color.Red, e.Elevated()).
		SetForeground(color.HiWhite)

	text := "%m"
	if s.Symbol != "" {
		text = s.Symbol
	}

	return prompt.NewSegmentPtr(" "+text+" ", *style), nil
}

func (t *flamesTheme) path(e *sys.Environment) (*prompt.Segment, error) {
	s, err := provider.NewPathStatus(e)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}

	style := color.NewStyle().
		SetBackground(color.Color(237)).
		SetForeground(color.HiWhite)

	text := s.Path
	if s.Symbol != "" {
		text = s.Symbol + " " + text
	}

	return prompt.NewSegmentPtr(" "+text+" ", *style), nil
}

func (t *flamesTheme) git(e *sys.Environment) (*prompt.Segment, error) {
	s, err := provider.NewGitStatus(e)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}

	style := color.NewStyle().
		SetForeground(color.HiWhite).
		SetBackground(color.Green).
		SetBackgroundIf(color.Yellow, s.Dirty)

	text := " " + s.Branch

	if s.StagedChanges {
		text += ""
	}

	return prompt.NewSegmentPtr(" "+text+" ", *style), nil
}

func (t *flamesTheme) kubernetes(e *sys.Environment) (*prompt.Segment, error) {
	s, err := provider.NewKubernetesStatus(e)
	if err != nil {
		return nil, err
	}
	if s == nil {
		return nil, nil
	}

	style := color.NewStyle().SetForeground(color.HiWhite)
	switch s.Env {
	case provider.KubernetesEnvDev:
		style.SetBackground(color.Green)
	case provider.KubernetesEnvCI:
		style.SetBackground(color.HiBlue)
	case provider.KubernetesEnvProd:
		style.SetBackground(color.Red)
	case provider.KubernetesEnvOther:
		style.SetBackground(color.Magenta)
	}

	var text string
	if s.Symbol != "" {
		text = s.Symbol
		if s.Identifier != "" {
			text += " " + s.Identifier
		}
		if s.Context != "" {
			text = s.Context + "@" + text
		}
	} else {
		if s.Context != "" {
			text = "☸ " + s.Context + "@" + s.Kubeconfig
		} else {
			text = "☸ " + s.Kubeconfig
		}
	}

	return prompt.NewSegmentPtr(" "+text+" ", *style), nil
}
