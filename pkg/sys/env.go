// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package sys

import "go.xrstf.de/promptomatic/pkg/config"

type Environment struct {
	EffectiveUserID  int
	WorkingDirectory string
	HomeDirectory    string
	GoPath           string
	Kubeconfig       string
	DumbTerminal     bool
	Config           *config.Config
}

func (e Environment) Elevated() bool {
	return e.EffectiveUserID == 0
}
