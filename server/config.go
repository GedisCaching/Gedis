package server

type Config struct {
	Address  string
	Password string
}

func DefaultConfig() *Config {
	return &Config{
		Address:  ":6379",
		Password: "",
	}
}
