package config

type Config struct {
	Rest      Rest      `yaml:"rest"`
	TGService TGService `yaml:"tg_service"`
}
