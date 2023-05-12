package config

import (
	"errors"
	"github.com/BurntSushi/toml"
	"io"
	"os"
)

func FromFile(path string, def interface{}) (interface{}, error) {
	file, err := os.Open(path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		return def, nil
	case err != nil:
		return nil, err
	}

	defer file.Close() //nolint:errcheck // The file is RO
	return FromReader(file, def)
}

// FromReader loads config from a reader instance.
func FromReader(reader io.Reader, def interface{}) (interface{}, error) {
	cfg := def
	_, err := toml.NewDecoder(reader).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
