package process

import (
	"path/filepath"
	"time"

	"github.com/Maycon-Santos/go-snake-backend/internal/projectpath"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Database struct {
	Driver       string        `mapstructure:"database_driver"`
	ConnURI      string        `mapstructure:"database_conn_uri"`
	MaxLifetime  time.Duration `mapstructure:"database_conn_max_lifetime"`
	IdleConns    uint          `mapstructure:"database_max_idle_conns"`
	MaxOpenConns uint          `mapstructure:"database_max_open_conns"`
}

type JWT struct {
	Secret           string        `mapstructure:"jwt_token_secret"`
	RefreshSecret    string        `mapstructure:"jwt_refresh_secret"`
	ExpiresIn        time.Duration `mapstructure:"jwt_expires_in"`
	RefreshExpiresIn time.Duration `mapstructure:"jwt_refresh_expires_in"`
}

type Env struct {
	AppName      string   `mapstructure:"app_name"`
	ServerPort   int      `mapstructure:"server_port"`
	RedisAddress string   `mapstructure:"redis_address"`
	JWT          JWT      `mapstructure:",squash"`
	Database     Database `mapstructure:",squash"`
	AllowOrigin  string   `mapstructure:"allow_origin"`
}

func NewEnv() (*Env, error) {
	env := &Env{}

	defaultEnvs, err := godotenv.Read(filepath.Join(projectpath.Root, ".env.sample"))
	if err != nil {
		return nil, err
	}

	for key, value := range defaultEnvs {
		viper.SetDefault(key, value)
	}

	viper.SetConfigType("env")
	viper.SetConfigFile(filepath.Join(projectpath.Root, ".env"))
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(env); err != nil {
		return nil, err
	}

	return env, nil
}
