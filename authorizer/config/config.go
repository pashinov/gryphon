package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Host string `yaml:"host" json:"host"`
		Port uint32 `yaml:"port" json:"port"`
	} `yaml:"server" json:"server"`
	Authorizer struct {
		Url string `yaml:"url" json:"url"`
	} `yaml:"authorizer" json:"authorizer"`
	IdentityServer struct {
		Url       string `yaml:"url" json:"url"`
		Resources struct {
			Identifier   string `yaml:"identifier" json:"identifier"`
			RefreshToken string `yaml:"refreshToken" json:"refreshToken"`
		} `yaml:"resources" json:"resources"`
	} `yaml:"identityServer" json:"identityServer"`
}

func NewConfig() *Config {
	return &Config{}
}

func (config *Config) Init(configFile string) error {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	return nil
}
