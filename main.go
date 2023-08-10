// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"go.xrstf.de/promptomatic/pkg/config"
	"go.xrstf.de/promptomatic/pkg/sys"
	"go.xrstf.de/promptomatic/pkg/theme"
	"go.xrstf.de/promptomatic/pkg/theme/flames"
	"go.xrstf.de/promptomatic/pkg/theme/simple"
	"go.xrstf.de/promptomatic/pkg/theme/zsh"

	"github.com/spf13/pflag"
)

// Project build specific vars
var (
	Tag    string
	Commit string
)

func printVersion() {
	fmt.Printf(
		"version: %s\nbuilt with: %s\ntag: %s\ncommit: %s\n",
		strings.TrimPrefix(Tag, "v"),
		runtime.Version(),
		Tag,
		Commit,
	)
}

type options struct {
	configFile string
	theme      string
	version    bool
}

func (o *options) AddPFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&o.configFile, "config", "c", "", "path to an optional config file (defaults to ~/.config/promptomatic.yaml)")
	fs.StringVarP(&o.theme, "theme", "t", "", "what theme to use to display the prompt (has priority over config.theme)")
	fs.BoolVarP(&o.version, "version", "V", o.version, "Show version info and exit immediately")
}

func main() {
	opt := options{}
	opt.AddPFlags(pflag.CommandLine)

	pflag.Parse()

	if opt.version {
		printVersion()
		return
	}

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
