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
	Debug        bool         `env:"DEBUG"         env-default:"false" yaml:"debug"`
	DBConfig     DBConfig     `env-required:"true" yaml:"db"`
	ServerConfig ServerConfig `env-required:"true" yaml:"server"`
}

type ServerConfig struct {
	Port        uint16        `env:"PORT"         env-default:"8080" yaml:"port"`
	Timeout     time.Duration `env:"TIMEOUT"      env-default:"5s"   yaml:"timeout"`
	IdleTimeout time.Duration `env:"IDLE_TIMEOUT" env-default:"30s"  yaml:"idle_timeout"`
}

type DBConfig struct {
	Host     string `env-required:"true" yaml:"host"`
	Port     uint16 `env-required:"true" yaml:"port"`
	User     string `env-required:"true" yaml:"user"`
	Password string `env:"DB_PASSWORD"   env-required:"true"`
	DBName   string `env-required:"true" yaml:"db_name"`
	SSLMode  string `env-required:"true" yaml:"ssl_mode"`
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
