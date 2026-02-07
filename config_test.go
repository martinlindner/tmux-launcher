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
	if cfg.AutoNewSession != true {
		t.Errorf("AutoNewSession = %v, want true", cfg.AutoNewSession)
	}
}

func TestConfigFlags(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		allowNested    bool
		autoAttach     bool
		autoNewSession bool
	}{
		{
			name:           "no flags",
			args:           nil,
			allowNested:    false,
			autoAttach:     true,
			autoNewSession: true,
		},
		{
			name:           "allow-nested",
			args:           []string{"--allow-nested"},
			allowNested:    true,
			autoAttach:     true,
			autoNewSession: true,
		},
		{
			name:           "no-auto-attach",
			args:           []string{"--no-auto-attach"},
			allowNested:    false,
			autoAttach:     false,
			autoNewSession: true,
		},
		{
			name:           "no-auto-new-session",
			args:           []string{"--no-auto-new-session"},
			allowNested:    false,
			autoAttach:     true,
			autoNewSession: false,
		},
		{
			name:           "all flags",
			args:           []string{"--allow-nested", "--no-auto-attach", "--no-auto-new-session"},
			allowNested:    true,
			autoAttach:     false,
			autoNewSession: false,
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
			if cfg.AutoNewSession != tt.autoNewSession {
				t.Errorf("AutoNewSession = %v, want %v", cfg.AutoNewSession, tt.autoNewSession)
			}
		})
	}
}

func TestConfigFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(cfgPath, []byte("allow_nested: true\nauto_attach: false\nauto_new_session: false\n"), 0644); err != nil {
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
	if cfg.AutoNewSession != false {
		t.Errorf("AutoNewSession = %v, want false", cfg.AutoNewSession)
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

	// Config file sets auto_new_session: true, flag overrides
	if err := os.WriteFile(cfgPath, []byte("auto_new_session: true\n"), 0644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	cfg, err = loadConfigFrom([]string{"--no-auto-new-session"}, cfgPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AutoNewSession != false {
		t.Errorf("AutoNewSession = %v, want false (flag should override file)", cfg.AutoNewSession)
	}
}
