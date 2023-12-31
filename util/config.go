package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config sotres all config of the aplicaton
// the values are read by viper from a config file
type Config struct {
	Environment          string        `mapstructure:"ENVIROMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MigrationURL         string        `mapstructure:"MIGRATION_URL"`
	HTTPServerAdress     string        `mapstructure:"HTTP_SERVER_ADRESS"`
	GRPCServerAdress     string        `mapstructure:"GRPC_SERVER_ADRESS"`
	TokenSymetricKey     string        `mapstructure:"TOKEN_SYMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

// Reads config file  or env var
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return

}
