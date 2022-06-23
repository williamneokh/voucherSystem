package config

import (
	"github.com/spf13/viper"
	"html/template"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ApiToken string `mapstructure:"API_TOKEN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}

type AppConfig struct {
	UserCache     bool
	TemplateCache map[string]*template.Template
	InProduction  bool
}
