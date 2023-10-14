// SPDX-FileCopyrightText: 2023 Christoph Mewes
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"go.xrstf.de/promptomatic/pkg/config"
	"go.xrstf.de/promptomatic/pkg/sys"
	"go.xrstf.de/promptomatic/pkg/theme"
	"go.xrstf.de/promptomatic/pkg/theme/flames"
	"go.xrstf.de/promptomatic/pkg/theme/simple"
	"go.xrstf.de/promptomatic/pkg/theme/zsh"

	"github.com/spf13/pflag"
)

// These variables get set by ldflags during compilation.
var (
	BuildTag    string
	BuildCommit string
	BuildDate   string // RFC3339 format ("2006-01-02T15:04:05Z07:00")
)

func printVersion() {
	fmt.Printf(
		"Promptomatic %s (%s), built with %s on %s\n",
		BuildTag,
		BuildCommit[:10],
		runtime.Version(),
		BuildDate,
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
	wd, err := os.Getwd()
	// An error can happen if the current directory of the terminal
	// is deleted by another process. In this case, Getwd will return
	// an error, but $PWD might still be set (though invalid in this
	// case, of course).
	if err != nil {
		wd = os.Getenv("PWD")
	}

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
