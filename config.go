package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type Config struct {
	AllowNested    bool `koanf:"allow_nested"`
	AutoAttach     bool `koanf:"auto_attach"`
	AutoNewSession bool `koanf:"auto_new_session"`
}

func loadConfig() (Config, error) {
	return loadConfigFrom(os.Args[1:], "")
}

func loadConfigFrom(args []string, cfgFile string) (Config, error) {
	k := koanf.New(".")

	// 1. Defaults
	k.Load(confmap.Provider(map[string]any{
		"allow_nested":     false,
		"auto_attach":      true,
		"auto_new_session": true,
	}, "."), nil)

	// 2. Config file (optional)
	if cfgFile == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			cfgFile = filepath.Join(home, ".config", "tmux-launcher", "config.yaml")
		}
	}
	if cfgFile != "" {
		if _, err := os.Stat(cfgFile); err == nil {
			if err := k.Load(file.Provider(cfgFile), yaml.Parser()); err != nil {
				return Config{}, err
			}
		}
	}

	// 3. CLI flags (override config file)
	f := pflag.NewFlagSet("tmux-launcher", pflag.ContinueOnError)
	f.Bool("allow-nested", false, "allow running inside an existing tmux session")
	f.Bool("no-auto-attach", false, "always show the TUI picker instead of auto-attaching")
	f.Bool("no-auto-new-session", false, "show the TUI picker even when no sessions exist")
	if err := f.Parse(args); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			os.Exit(0)
		}
		return Config{}, err
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return Config{}, err
	}

	// Apply flag overrides (handle name mismatch and negation)
	if f.Changed("allow-nested") {
		v, _ := f.GetBool("allow-nested")
		cfg.AllowNested = v
	}
	if f.Changed("no-auto-attach") {
		v, _ := f.GetBool("no-auto-attach")
		cfg.AutoAttach = !v
	}
	if f.Changed("no-auto-new-session") {
		v, _ := f.GetBool("no-auto-new-session")
		cfg.AutoNewSession = !v
	}

	return cfg, nil
}
