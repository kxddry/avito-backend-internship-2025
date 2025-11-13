package startup

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

func ReadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	flag.StringVar(&configPath, "config", configPath, "Path to configuration file")

	flag.Parse()
	var c Config
	if err := cleanenv.ReadConfig(configPath, &c); err != nil {
		return nil, err
	}
	if err := c.Validate(); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}
	return &c, nil
}

func (c *Config) Validate() error {
	return nil // при изменении конфига можно реализовать, пока что ничего проверять не надо
}

type Config struct {
	Debug        bool         `yaml:"debug" env:"DEBUG" env-default:"false"`
	DBConfig     DBConfig     `yaml:"db" env-required:"true"`
	ServerConfig ServerConfig `yaml:"server" env-required:"true"`
}

type ServerConfig struct {
	Port        uint16        `yaml:"port" env:"PORT" env-default:"8080"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"IDLE_TIMEOUT" env-default:"30s"`
}

type DBConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     uint16 `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
	SSLMode  string `yaml:"ssl_mode" env-required:"true"`
}

func (d *DBConfig) DSN() string {
	userInfo := url.UserPassword(d.User, d.Password)
	u := &url.URL{
		Scheme: "postgresql",
		User:   userInfo,
		Host:   fmt.Sprintf("%s:%d", d.Host, d.Port),
		Path:   d.DBName,
	}

	query := u.Query()
	query.Set("sslmode", d.SSLMode)
	u.RawQuery = query.Encode()

	return u.String()
}
