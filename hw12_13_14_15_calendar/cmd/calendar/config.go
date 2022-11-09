package main

import (
	"fmt"
	"github.com/spf13/viper"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf `json:"logger"`
	DB     DBConf     `json:"db"`
	HTTP   HTTPConf   `json:"http"`
}

type LoggerConf struct {
	Level string `json:"level"`
	File  string `json:"file"`
}

type DBConf struct {
	Method           string `json:"method"`
	ConnectionString string `json:"connection_string"`
}

type HTTPConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	GrpcPort string `json:"grpc_port"`
}

func NewConfig() (Config, error) {
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		return Config{}, fmt.Errorf("fatal error config file: %w", err)
	}

	return Config{
		LoggerConf{Level: viper.GetString("logger.level"), File: viper.GetString("logger.file")},
		DBConf{Method: viper.GetString("db.method"), ConnectionString: viper.GetString("db.connection_string")},
		HTTPConf{Host: viper.GetString("http.host"), Port: viper.GetString("http.port"), GrpcPort: viper.GetString("http.grpc_port")},
	}, nil
}
