package config

import "crypto/rsa"

type Config struct {
	Db  Db  `mapstructure:"db"`
	App App `mapstructure:"app"`
}

type Db struct {
	User   string `mapstructure:"user"`
	Passwd string `mapstructure:"passwd"`
	Addr   string `mapstructure:"addr"`
	Port   string `mapstructure:"port"`
	Name   string `mapstructure:"name"`
}

type App struct {
	CertificateKey     string          `mapstructure:"certificate_key"`
	CertificateKeyPath string          `mapstructure:"certificate_key_path"`
	PrivateKey         *rsa.PrivateKey `mapstructure:"-"`
}
