package process

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Env struct {
	AppName    string `mapstructure:"app_name"`
	ServerPort int    `mapstructure:"server_port"`
}

func NewEnv() (*Env, error) {
	env := &Env{}

	defaultEnvs, err := godotenv.Read(".env.sample")
	if err != nil {
		return nil, err
	}

	for key, value := range defaultEnvs {
		viper.SetDefault(key, value)
	}

	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	if err := viper.Unmarshal(env); err != nil {
		return nil, err
	}

	return env, nil
}
