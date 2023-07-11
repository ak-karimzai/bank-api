package util

import (
	"time"

	"github.com/spf13/viper"
)

const (
	DevelopmentEnvironment = "development"
	ProductionEnvironment  = "release"
)

type Config struct {
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmtricKey     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	MigrationDir         string        `mapstructure:"MIGRATION_DIR"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	Environment          string        `mapstructure:"ENVIRONMENT"`
	RedisServerAddress   string        `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(path + "/.env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
