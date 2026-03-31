package domain

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

var globalConfig *Config

type Database struct {
	User         string `json:"user"`
	Pass         string `json:"pass"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DBName       string `json:"db_name"`
	ConnLifeTime int    `json:"connection_life_time"`
	MaxConnPool  int    `json:"pool"`
	MaxIdleConn  int    `json:"idle_connections"`
}

type Redis struct {
	Host        string `json:"host"`
	Pass        string `json:"pass"`
	Database    int    `json:"database"`
	Pool        int    `json:"pool"`
	Port        int    `json:"port"`
	MinIdleConn int    `json:"idle_connections"`
}

type Config struct {
	Database Database `json:"database"`
	Redis    Redis    `json:"redis"`

	PrivacyServiceHost string `json:"privacy_service_host"`
	ListenPort         int32    `json:"listen_port"`
	Http2ListenPort    int    `json:"http2_listen_port"`

	SysLogAddr   string `json:"syslog_address"`
	LogsToStdout bool   `json:"logs_to_stdout"`
	GinLogs      bool   `json:"gin_logs"`
	ErrLevel     int    `json:"error_level"`
	JWTSecret    string `json:"jwt_secret"`
	JWTExpiry    int64  `json:"jwt_expiry"`
}

func New(configFilePath string) *Config {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal().Msgf("can't read config file: %v", err)
	}

	err = json.Unmarshal(file, &globalConfig)
	if err != nil {
		log.Fatal().Msgf("can't unmarshal config: %v", err)
	}

	return globalConfig
}
