package cfg

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Host        string `mapstructure:"HOST"`
	Port        string `mapstructure:"PORT"`
	DBUsername  string `mapstructure:"MYSQL_USER"`
	DBPassword  string `mapstructure:"MYSQL_PASS"`
	DBHost      string `mapstructure:"MYSQL_HOST"`
	DBPort      string `mapstructure:"MYSQL_PORT"`
	DBName      string `mapstructure:"MYSQL_DB"`
	DBNameTest  string `mapstructure:"MYSQL_DB_TEST"`
	DBUrl       string
}

func LoadConfig(name string, path string) (config Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("config: %v", err)
		return
	}
	return
}
