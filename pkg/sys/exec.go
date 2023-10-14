// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package sys

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func RunCommand(cwd string, cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)

	command.Stdout = &stdout
	command.Stderr = &stderr

	if cwd != "" {
		command.Dir = cwd
	}

	if err := command.Run(); err != nil {
		return "", errors.New(strings.TrimSpace(stderr.String()))
	}

	return strings.TrimSpace(stdout.String()), nil
}
