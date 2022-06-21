package config

import "github.com/spf13/viper"

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ApiToken string `mapstructure:"API_TOKEN"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	//viper.AddConfigPath("/app/voucherApp")
	//viper.AddConfigPath("/app/voucherapi")
	//viper.AddConfigPath("voucherapi")
	viper.AddConfigPath("app")
	//viper.AddConfigPath(".")
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
