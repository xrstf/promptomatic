// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package provider

import "go.xrstf.de/promptomatic/pkg/sys"

type WhoamiStatus struct {
	Symbol string
}

func NewWhoamiStatus(e *sys.Environment) (*WhoamiStatus, error) {
	s := &WhoamiStatus{
		Symbol: e.Config.HostSymbol,
	}

	return s, nil
}
