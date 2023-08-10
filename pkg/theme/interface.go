// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package theme

import "go.xrstf.de/promptomatic/pkg/sys"

type Theme interface {
	Render(e *sys.Environment) string
}
