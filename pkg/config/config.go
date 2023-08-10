// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package config

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Theme      string `yaml:"theme"`
	HostSymbol string `yaml:"hostSymbol"`
}

func Load(filename string) (*Config, error) {
	if filename == "" {
		filename = filepath.Join(os.Getenv("HOME"), ".config", "promptomatic.yaml")
		if info, err := os.Stat(filename); err != nil || info.Size() == 0 {
			return &Config{}, nil
		}
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	c := &Config{}
	if err := yaml.NewDecoder(f).Decode(c); err != nil {
		// gracefully handle files with only whitespace
		if errors.Is(err, io.EOF) {
			return &Config{}, nil
		}

		return nil, err
	}

	return c, nil
}

func Validate(c *Config) error {
	return nil
}
