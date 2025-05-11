package config

import (
	"context"
	"encoding/json"
	"os"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v2"
)

// LoadConfig loads configuration from a file.
// It also supports loading configuration from vault.
// If vaultConfig is set, it will load the configuration from vault.
// If vaultConfig is not set, it will load the entire configuration from the file.
// NOTE: If vaultConfig is set, vaultOptions must be provided.
func LoadConfig(
	ctx context.Context,
	cfg any,
	opts ...Option,
) (err error) {
	// create default config options.
	opt := newDefaultConfigOptions()
	for _, o := range opts {
		o.apply(opt)
	}

	// load config from file.
	jsonFile, err := os.Open(opt.cfgPath)
	if err != nil {
		return err
	}

	// parse json file.
	err = json.NewDecoder(jsonFile).Decode(&cfg)
	if err != nil {
		return err
	}

	if !opt.disableDefaults {
		if err := defaults.Set(cfg); err != nil {
			return err
		}
	}

	// validate config.
	err = validator.New().Struct(cfg)
	if err != nil {
		return err
	}

	return err
}

// LoadFromYAML loads configuration from a yaml file.
func LoadFromYAML(cfg any, path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(file, cfg); err != nil {
		return err
	}

	return nil
}
