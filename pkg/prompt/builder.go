// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package prompt

import (
	"go.xrstf.de/promptomatic/pkg/color"
)

type Builder struct {
	transition styleTransitioner
	segments   []Segment
}

func NewBuilder(t styleTransitioner) *Builder {
	return &Builder{
		transition: t,
		segments:   []Segment{},
	}
}

func (b *Builder) AddSegment(s Segment) *Builder {
	b.segments = append(b.segments, s)
	return b
}

func (b *Builder) Render() string {
	var (
		result   string
		curStyle color.Style
	)

	for _, segment := range b.segments {
		if result != "" {
			result += b.transition(curStyle, segment.style, false)
		}

		result += segment.style.Apply(segment.text)
		curStyle = segment.style.DeepCopy()
	}

	// add one final transition from the promptomatic into the terminal background
	if result != "" {
		result += b.transition(curStyle, color.EmptyStyle, true)
	}

	return result + " "
}

type Segment struct {
	text  string
	style color.Style
}

func NewSegment(text string, style color.Style) Segment {
	return *NewSegmentPtr(text, style)
}

func NewSegmentPtr(text string, style color.Style) *Segment {
	return &Segment{
		text:  text,
		style: style,
	}
}
