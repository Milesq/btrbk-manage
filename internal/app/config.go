package app

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"
	"milesq.dev/btrbk-manage/internal/btrbk"
)

const DefaultConfigPath = "/etc/btrbk-manage/config.yaml"

type Config struct {
	configPath                string
	BtrbkConfigFile           string   `yaml:"btrbk_config_file"`
	DefaultSubvolsRestoreList []string `yaml:"default_subvols_restore_list"`
	OldFormat                 string   `yaml:"old_format"`
	Paths                     Paths    `yaml:"paths"`
}

type Paths struct {
	Snaps     string `yaml:"snaps"`
	Target    string `yaml:"target"`
	Meta      string `yaml:"meta,omitempty"`
	MetaTrash string `yaml:"meta_trash,omitempty"`
	Hooks     string `yaml:"hooks,omitempty"`
}

func (p Paths) String() string {
	return fmt.Sprintf("Paths{Snaps: %q, Target: %q, Meta: %q, MetaTrash: %q}",
		p.Snaps, p.Target, p.Meta, p.MetaTrash)
}

func (c *Config) String() string {
	return fmt.Sprintf("Config{BtrbkConfigFile: %q, DefaultSubvolsRestoreList: %v, OldFormat: %q, Paths: %s}",
		c.BtrbkConfigFile, c.DefaultSubvolsRestoreList, c.OldFormat, c.Paths)
}

func LoadConfig(configPath, project string) (*Config, error) {
	actualConfigPath := DefaultConfigPath
	if project != "" {
		actualConfigPath = fmt.Sprintf("/etc/btrbk-manage/config.%s.yaml", project)
	} else if configPath != "" {
		actualConfigPath = configPath
	}

	cfg, err := readConfig(actualConfigPath)
	if err != nil {
		return nil, err
	}

	if err := cfg.detectMissing(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func readConfig(path string) (*Config, error) {
	var cfg Config
	cfg.configPath = path
	data, err := os.ReadFile(path)

	if err != nil {
		if path == DefaultConfigPath {
			return &cfg, nil
		}

		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &cfg, nil
}

func (c *Config) detectMissing() error {
	if c.Paths.Snaps == "" || c.Paths.Target == "" {
		results, err := btrbk.List(c.BtrbkConfigFile)

		if err != nil {
			return err
		}

		if c.Paths.Snaps == "" {
			c.Paths.Snaps = results[0].SnapPath
		}

		if c.Paths.Target == "" {
			c.Paths.Target = filepath.Dir(results[0].Source)
		}
	}

	if c.OldFormat == "" {
		c.OldFormat = "{{.SubvolName}}.old"
	}

	if c.Paths.Meta == "" {
		c.Paths.Meta = c.Paths.Snaps + "/.meta"
	}

	if c.Paths.MetaTrash == "" {
		c.Paths.MetaTrash = c.Paths.Snaps + "/.meta/.trash"
	}

	if c.Paths.Hooks == "" {
		c.Paths.Hooks = path.Dir(c.configPath) + "/hooks"
	}

	return nil
}
