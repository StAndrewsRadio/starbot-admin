package cfg

import (
	"github.com/pelletier/go-toml"
)

type Config struct {
	tree *toml.Tree
}

// Loads a new config struct.
func New(path string) (*Config, error) {
	tree, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	return &Config{tree: tree}, nil
}

// Gets a string from the config, given a key.
func (config *Config) GetString(key string) string {
	return config.tree.Get(key).(string)
}

// Gets a slice of strings from the config, given a key.
func (config *Config) GetStrings(key string) []string {
	interfaces := config.tree.Get(key).([]interface{})
	strings := make([]string, len(interfaces))

	for i := range interfaces {
		strings[i] = interfaces[i].(string)
	}

	return strings
}

// Gets a boolean from the config, given a key.
func (config *Config) GetBool(key string) bool {
	return config.tree.Get(key).(bool)
}
