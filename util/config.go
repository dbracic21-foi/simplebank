package util

import "github.com/spf13/viper"

// Config sotres all config of the aplicaton
// the values are read by viper from a config file
type Config struct {
	DBDRIVER     string `mapstructure:"DB_DRIVER"`
	DBSOURCE     string `mapstructure:"DBS_SOURCE"`
	ServerAdress string `mapstructure:"SERVER_ADRESS"`
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
