package storage

type Config struct {
	DatabaseURI string `toml:"database_uri"`
}

func NewConfig() *Config {
	return &Config{}
}
