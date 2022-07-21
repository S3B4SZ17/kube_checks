package healthcheck

import (
	"os"

	"gopkg.in/yaml.v2"
)

// NewConfig returns a new decoded Config struct
func ReadYaml(configPath string) (*Config, error) {
    // Create config structure
    config := &Config{}

    // Open config file
    file, err := os.Open(configPath)
    if err != nil {
      return nil, err
    }
    defer file.Close()

    // Init new YAML decode
    d := yaml.NewDecoder(file)

    // Start YAML decoding from file
    if err := d.Decode(&config); err != nil {
        return nil, err
    }

    return config, nil
}

type Config struct {
  Vault_url string `yaml:"vault_url"`
  Endpoints []struct {
    App string `yaml:"app"`
    Url string `yaml:"url"`
  } `yaml:"endpoints"`
}