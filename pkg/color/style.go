// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package color

import "fmt"

type Style struct {
	Foreground *Color
	Background *Color
	Bold       bool
}

var EmptyStyle = Style{}

func NewStyle() *Style {
	return &Style{}
}

func (s *Style) SetForeground(c Color) *Style {
	s.Foreground = &c
	return s
}

func (s *Style) SetForegroundIf(c Color, condition bool) *Style {
	if condition {
		return s.SetForeground(c)
	}

	return s
}

func (s *Style) SetBackground(c Color) *Style {
	s.Background = &c
	return s
}

func (s *Style) SetBackgroundIf(c Color, condition bool) *Style {
	if condition {
		return s.SetBackground(c)
	}

	return s
}

func (s *Style) SetBold(bold bool) *Style {
	s.Bold = bold
	return s
}

func (s Style) DeepCopy() Style {
	return Style{
		Foreground: s.Foreground,
		Background: s.Background,
		Bold:       s.Bold,
	}
}

func (s Style) Apply(text string) string {
	if s.Foreground != nil {
		text = fmt.Sprintf("%%F{%d}%s%%f", *s.Foreground, text)
	}

	if s.Background != nil {
		text = fmt.Sprintf("%%K{%d}%s%%k", *s.Background, text)
	}

	if s.Bold {
		text = "%B" + text + "%b"
	}

	return text
}
