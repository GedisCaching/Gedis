package redis

type Config struct {
	Address  string
	Password string
}

func DefaultConfig() *Config {
	return &Config{
		Address:  "localhost:6379",
		Password: "",
	}
}
