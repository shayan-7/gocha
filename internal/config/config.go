package config

import "github.com/spf13/viper"

type Config struct {
	ListenAddr string `mapstructure:"listen_addr"`
	RedisAddr  string `mapstructure:"redis_addr"`
	RedisChan  string `mapstructure:"redis_chan"`
	DBUsername string `mapstructure:"db_username"`
	DBPassword string `mapstructure:"db_password"`
	DBAddr     string `mapstructure:"db_addr"`
	DBName     string `mapstructure:"db_name"`
}

func LoadConfig(path string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("gocha")
	viper.SetConfigType("yml")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
