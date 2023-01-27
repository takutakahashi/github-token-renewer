package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL            string         `yaml:"url"`
	PrivateKeyPath string         `yaml:"privateKeyPath"`
	AppID          int64          `yaml:"appID"`
	Installations  []Installation `yanl:"installations"`
}

type Installation struct {
	ID     int64  `yaml:"id"`
	Output Output `yaml:"output"`
}

type Output struct {
	KubernetesSecret *OutputKubernetesSecret `yaml:"kubernetesSecret"`
}

type OutputKubernetesSecret struct {
	SecretName      string `yaml:"secretName"`
	SecretNamespace string `yaml:"secretNamespace"`
	Key             string `yaml:"key"`
}

func Load(p string) (*Config, error) {
	ret := &Config{}
	buf, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(buf, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func GetInstallation(cfg *Config, id int64) *Installation {
	if cfg.Installations == nil {
		return nil
	}
	for _, i := range cfg.Installations {
		if i.ID == id {
			return &i
		}
	}
	return nil
}
