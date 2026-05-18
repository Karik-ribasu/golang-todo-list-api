package config

import (
	"fmt"

	viper "github.com/spf13/viper"
)

func readMergedConfig(v *viper.Viper) (cfg Config, err error) {
	if err = v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("unmarshal config: %w", err)
	}
	return cfg, nil
}

func ReadConfig() (cfg Config, err error) {
	v := viper.New()
	v.AddConfigPath("./")
	v.SetConfigName("config")
	v.SetConfigType("toml")
	if err = v.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("read config: %w", err)
	}
	return readMergedConfig(v)
}
