package cfg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestEnv(t *testing.T, cfgContent string) string {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Write the config content to a temporary file
	cfgFilePath := filepath.Join(tempDir, "config.env")
	err := os.WriteFile(cfgFilePath, []byte(cfgContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Set the CFG_PATH environment variable
	os.Setenv("CFG_PATH", tempDir)

	return tempDir
}

func TestLoadConfig_Success(t *testing.T) {
	// Setup test environment
	cfgContent := `
ENVIRONMENT=production
HOST=localhost
PORT=8080
MYSQL_USER=root
MYSQL_PASS=password
MYSQL_HOST=127.0.0.1
MYSQL_PORT=3306
MYSQL_DB=mydb
MYSQL_DB_TEST=mydb_test
APP_DOMAIN=example.com
`
	setupTestEnv(t, cfgContent)

	// Load config
	config, err := LoadConfig("config")
	assert.NoError(t, err)

	// Validate config
	assert.Equal(t, "production", config.Environment)
	assert.Equal(t, "localhost", config.Host)
	assert.Equal(t, "8080", config.Port)
	assert.Equal(t, "root", config.DBUsername)
	assert.Equal(t, "password", config.DBPassword)
	assert.Equal(t, "127.0.0.1", config.DBHost)
	assert.Equal(t, "3306", config.DBPort)
	assert.Equal(t, "mydb", config.DBName)
	assert.Equal(t, "mydb_test", config.DBNameTest)
	assert.Equal(t, "example.com", config.AppDomain)
	assert.Equal(t, "root:password@tcp(127.0.0.1:3306)/mydb?parseTime=true", config.DBUrl)
}

func TestLoadConfig_MissingCFGPath(t *testing.T) {
	// Unset CFG_PATH
	os.Unsetenv("CFG_PATH")

	// Load config
	_, err := LoadConfig("config")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "CFG_PATH environment variable is not set")
}

func TestLoadConfig_InvalidConfigFile(t *testing.T) {
	// Setup test environment with an empty config file
	setupTestEnv(t, "")

	// Load config
	_, err := LoadConfig("config")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config file is empty") // Updated to match the actual error message
}

func TestConstructDBUrl_TestEnvironment(t *testing.T) {
	config := Config{
		Environment: "test",
		DBUsername:  "root",
		DBPassword:  "password",
		DBHost:      "127.0.0.1",
		DBPort:      "3306",
		DBName:      "mydb",
		DBNameTest:  "mydb_test",
	}

	dbUrl := constructDBUrl(config)
	assert.Equal(t, "root:password@tcp(127.0.0.1:3306)/mydb_test?parseTime=true", dbUrl)
}

func TestConstructDBUrl_ProductionEnvironment(t *testing.T) {
	config := Config{
		Environment: "production",
		DBUsername:  "root",
		DBPassword:  "password",
		DBHost:      "127.0.0.1",
		DBPort:      "3306",
		DBName:      "mydb",
		DBNameTest:  "mydb_test",
	}

	dbUrl := constructDBUrl(config)
	assert.Equal(t, "root:password@tcp(127.0.0.1:3306)/mydb?parseTime=true", dbUrl)
}
