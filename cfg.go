package cfg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	AppDomain   string `mapstructure:"APP_DOMAIN"`
	DBUrl       string
}

// LoadConfig loads the configuration from the specified file and path.
// The path is read from the CFG_PATH environment variable.

func LoadConfig(name string) (Config, error) {
	// Read the config path from the CFG_PATH environment variable
	cfgPath := os.Getenv("CFG_PATH")
	if cfgPath == "" {
		return Config{}, fmt.Errorf("CFG_PATH environment variable is not set")
	}

	viper.AddConfigPath(cfgPath)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	// Check if the config file exists and is not empty
	cfgFile := filepath.Join(cfgPath, name+".env")
	fileInfo, err := os.Stat(cfgFile)
	if err != nil {
		return Config{}, fmt.Errorf("config file not found: %w", err)
	}
	if fileInfo.Size() == 0 {
		return Config{}, fmt.Errorf("config file is empty")
	}

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Construct the DBUrl based on the environment
	config.DBUrl = constructDBUrl(config)

	return config, nil
}

// constructDBUrl constructs the database URL based on the configuration.
func constructDBUrl(config Config) string {
	dbName := config.DBName
	if config.Environment == "test" {
		dbName = config.DBNameTest
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		dbName,
	)
}

// MustLoadConfig is a convenience function that logs and exits on error.
func MustLoadConfig(name string) Config {
	config, err := LoadConfig(name)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return config
}
