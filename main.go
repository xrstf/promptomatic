// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"
	"os"

	"go.xrstf.de/promptomatic/pkg/config"
	"go.xrstf.de/promptomatic/pkg/sys"
	"go.xrstf.de/promptomatic/pkg/theme"
	"go.xrstf.de/promptomatic/pkg/theme/flames"
	"go.xrstf.de/promptomatic/pkg/theme/simple"
	"go.xrstf.de/promptomatic/pkg/theme/zsh"

	"github.com/spf13/pflag"
)

type options struct {
	configFile string
	theme      string
}

func (o *options) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.configFile, "config", "c", "", "path to an optional config file")
	fs.StringVarP(&o.theme, "theme", "t", "", "what theme to use to display the promptomatic (has priority overr config.theme)")
}

func main() {
	opt := options{}
	opt.AddPFlags(pflag.CommandLine)

	pflag.Parse()

	cfg, err := config.Load(opt.configFile)
	if err != nil {
		log.Fatalf("Failed to load config file %q: %v", opt.configFile, err)
	}

	var t theme.Theme

	effectiveTheme := opt.theme
	if effectiveTheme == "" {
		effectiveTheme = cfg.Theme
		if effectiveTheme == "" {
			effectiveTheme = "auto"
		}
	}

	switch effectiveTheme {
	case "zsh":
		t = zsh.New()
	case "simple":
		t = simple.New()
	case "flames":
		t = flames.New()
	case "auto":
		t = autodetectTheme()
	default:
		log.Printf("Unknown theme %q. Falling back to 'simple'.", opt.theme)
		t = simple.New()
	}

	env := createEnvironment(cfg)

	fmt.Print(t.Render(env))
}

func autodetectTheme() theme.Theme {
	term := os.Getenv("TERM")

	switch term {
	case "screen-256color":
		fallthrough
	case "xterm-256color":
		return flames.New()
	case "xterm":
		fallthrough
	default:
		return simple.New()
	}
}

func createEnvironment(cfg *config.Config) *sys.Environment {
	wd, _ := os.Getwd()

	env := &sys.Environment{
		EffectiveUserID:  os.Geteuid(),
		GoPath:           os.Getenv("GOPATH"),
		HomeDirectory:    os.Getenv("HOME"),
		Kubeconfig:       os.Getenv("KUBECONFIG"),
		WorkingDirectory: wd,
		Config:           cfg,
	}

	return env
}
