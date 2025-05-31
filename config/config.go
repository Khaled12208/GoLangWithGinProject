package config

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Server   ServerConfig    `mapstructure:"server"`
	Database DatabaseConfig  `mapstructure:"database"`
	Sharding ShardingConfig `mapstructure:"sharding"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port string `mapstructure:"port"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

type ShardingConfig struct {
	Enabled bool            `mapstructure:"enabled"`
	Shards  []DatabaseConfig `mapstructure:"shards"`
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type LoggerConfig struct {
	Level string
	File  string
}

// Load reads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	
	// Set defaults
	viper.SetDefault("server.port", "8888")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("sharding.enabled", false)
	
	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvPrefix("APP")
	
	// Map environment variables
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.username", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.dbname", "DB_NAME")
	viper.BindEnv("sharding.enabled", "DB_SHARDING_ENABLED")
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	
	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	
	return &config, nil
} 