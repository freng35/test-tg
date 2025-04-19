package config

type Config struct {
	GRPC     GRPC     `yaml:"grpc"`
	Postgres Postgres `yaml:"postgres"`
}
