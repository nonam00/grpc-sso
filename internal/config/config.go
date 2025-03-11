package config

import (
	"flag"
	"os"
	"time"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    Env         string        `yaml:"env" env-default:"local"`
    TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
    PGConn      PGConn        `yaml:"postgres_connection" env-required:"./data"`
    GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
    Port    int           `yaml:"port"` 
    Timeout time.Duration `yaml:"timeout"`
}

type PGConn struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
}

func MustLoad() *Config {
    path := fetchConfigPath()
    if path == "" {
        panic("config path is empty")
    }

    return MustLoadByPath(path);
}

func MustLoadByPath(configPath string) *Config {
    if _, err := os.Stat(configPath); os.IsNotExist(err) {
        panic("config file does not exist: " + configPath)
    }

    var cfg Config

    if err:=cleanenv.ReadConfig(configPath, &cfg); err != nil {
        panic("failed to read config: " + err.Error())
    }

    return &cfg;
}

// fetchConfigPath fetches config path from command line flagg or env variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
    var res string

    // --config="path/to/config.yaml"
    flag.StringVar(&res, "config", "", "path to config file")
    flag.Parse()
    
    if res == "" {
        res = os.Getenv("CONFIG_PATH")
    }
    
    return res
}
