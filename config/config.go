package config

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/larspensjo/config"
)

const (
	// SectionDefault is the config file default section.
	sectionDefault = "Default"
)

// Config contains config info.
type Config struct {
	cfg *config.Config
}

// LoadConfig load the config file.
func LoadConfig(file string) (*Config, error) {
	c := new(Config)

	c.cfg, err = config.ReadDefault(file)
	if err != nil {
		return nil, errors.Errorf("failed to read file, %v", file)
	}

	return c, nil
}

// GetBool get bool value from config file with section and key.
func (c *Config) GetBool(section, key string) bool {
	if c.cfg.HasSection(section) {
		val, err := c.config.Bool(section, key)
		if err != nil {
			panic(fmt.Sprintf("get bool failure, section:%v key:%v err:%v", section, key, err))
		}
		return val
	}

	panicSection(section)
	return false
}

// GetInt get int value from config file with section and key.
func (cfg *Config) GetInt(section, key string) int {
	if cfg.config.HasSection(section) {
		val, err := cfg.config.Int(section, key)
		if err != nil {
			panic(fmt.Sprintf("get value failure, section:%v key:%v err:%v", section, key, err))
		}
		return val
	}

	panicSection(section)
	return 0
}

// GetString get string value from config file with section and key.
func (cfg *Config) GetString(section, key string) string {
	if cfg.config.HasSection(section) {
		val, err := cfg.config.String(section, key)
		if err != nil {
			panic(fmt.Sprintf("get string failure, section:%v key:%v err:%v", section, key, err))
		}
		return val
	}

	panicSection(section)
	return ""
}

func panicSection(section string) {
	panic(fmt.Sprintf("the section: %v not exist", section))
}
