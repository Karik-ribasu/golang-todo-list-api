package config

import (
	viper "github.com/spf13/viper"
)

func ReadConfig() (cfg Config, err error) {
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.ReadInConfig()

	err = viper.Unmarshal(&cfg)
	return cfg, err
}
