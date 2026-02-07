package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	cfg, err := loadConfigFrom(nil, "/nonexistent/config.yaml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AllowNested != false {
		t.Errorf("AllowNested = %v, want false", cfg.AllowNested)
	}
	if cfg.AutoAttach != true {
		t.Errorf("AutoAttach = %v, want true", cfg.AutoAttach)
	}
}

func TestConfigFlags(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		allowNested bool
		autoAttach  bool
	}{
		{
			name:        "no flags",
			args:        nil,
			allowNested: false,
			autoAttach:  true,
		},
		{
			name:        "allow-nested",
			args:        []string{"--allow-nested"},
			allowNested: true,
			autoAttach:  true,
		},
		{
			name:        "no-auto-attach",
			args:        []string{"--no-auto-attach"},
			allowNested: false,
			autoAttach:  false,
		},
		{
			name:        "both flags",
			args:        []string{"--allow-nested", "--no-auto-attach"},
			allowNested: true,
			autoAttach:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := loadConfigFrom(tt.args, "/nonexistent/config.yaml")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if cfg.AllowNested != tt.allowNested {
				t.Errorf("AllowNested = %v, want %v", cfg.AllowNested, tt.allowNested)
			}
			if cfg.AutoAttach != tt.autoAttach {
				t.Errorf("AutoAttach = %v, want %v", cfg.AutoAttach, tt.autoAttach)
			}
		})
	}
}

func TestConfigFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(cfgPath, []byte("allow_nested: true\nauto_attach: false\n"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err := loadConfigFrom(nil, cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AllowNested != true {
		t.Errorf("AllowNested = %v, want true", cfg.AllowNested)
	}
	if cfg.AutoAttach != false {
		t.Errorf("AutoAttach = %v, want false", cfg.AutoAttach)
	}
}

func TestConfigFlagOverridesFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	// Config file sets auto_attach: false
	if err := os.WriteFile(cfgPath, []byte("auto_attach: false\n"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	// But no flag is passed, so file value should win
	cfg, err := loadConfigFrom(nil, cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AutoAttach != false {
		t.Errorf("AutoAttach = %v, want false (from file)", cfg.AutoAttach)
	}

	// Config file sets allow_nested: true, flag overrides with --allow-nested=false
	if err := os.WriteFile(cfgPath, []byte("allow_nested: true\n"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err = loadConfigFrom([]string{"--allow-nested=false"}, cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AllowNested != false {
		t.Errorf("AllowNested = %v, want false (flag should override file)", cfg.AllowNested)
	}
}
