package cfg

type Config struct {
	LogLevel   int
	ListenPort int
}

func NewConfig() Config {
	return Config{
		LogLevel:   3,
		ListenPort: 8080,
	}
}
