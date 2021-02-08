package config

import "github.com/spf13/viper"

type Config struct {
	Server struct {
		Host string `yaml:"host" json:"host"`
		Port uint32 `yaml:"port" json:"port"`
	} `yaml:"server" json:"server"`
	IdentityServer struct {
		Url       string `yaml:"url" json:"url"`
		Resources struct {
			UserInfo string `yaml:"userinfo" json:"userinfo"`
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
