// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package prompt

import "go.xrstf.de/promptomatic/pkg/color"

type styleTransitioner func(old, new color.Style, isEnd bool) string

func NewSymbolTransition(symbol string) styleTransitioner {
	return func(old, new color.Style, _ bool) string {
		merged := color.NewStyle()

		if back := old.Background; back != nil {
			merged.SetForeground(*back)
		}

		if back := new.Background; back != nil {
			merged.SetBackground(*back)
		}

		return merged.Apply(symbol)
	}
}

func NewEmptyTransition() styleTransitioner {
	return func(old, new color.Style, isEnd bool) string {
		if isEnd {
			return ""
		}

		return " "
	}
}
