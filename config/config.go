package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

type GlobalConfig struct {
	Interval  int          `hcl:"interval"`
	Timeout   int          `hcl:"timeout"`
	Threshold int          `hcl:"threshold"`
	Log       *Log         `hcl:"log,block"`
	Alert     *AlertConfig `hcl:"alert,block"`
}

type LogMode string

type Log struct {
	Enabled bool    `hcl:"enabled"`
	Path    string  `hcl:"path"`
	Format  string  `hcl:"format"`
	MaxSize int     `hcl:"maxsize"`
	LogMode LogMode `hcl:"mode"`
}

type AlertConfig struct {
	Telegram *TelegramAlert `hcl:"telegram,block"`
	Email    *EmailAlert    `hcl:"email,block"`
}

type EmailAlert struct {
	Enabled bool                  `hcl:"enabled"`
	Config  *EmailConf            `hcl:"config,block"`
	With    *AlertConfigEmailWith `hcl:"with,block"`
}

type EmailConf struct {
	Host     string `hcl:"host"`
	Port     int    `hcl:"port"`
	Username string `hcl:"username"`
	Password string `hcl:"password"`
	From     string `hcl:"from"`
	To       string `hcl:"to"`
	// Subject  string  `hcl:"subject"`
	// Message  string  `hcl:"message"`
	CC *string `hcl:"cc"`
}

type AlertConfigEmailWith struct {
	Extends  string  `hcl:"extends,label"`
	Host     *string `hcl:"host"`
	Port     *int    `hcl:"port"`
	Username *string `hcl:"username"`
	Password *string `hcl:"password"`
	From     *string `hcl:"from"`
	To       *string `hcl:"to"`
	// Subject  *string `hcl:"subject"`
	// Message  *string `hcl:"message"`
	CC *string `hcl:"cc"`
}

type TelegramAlert struct {
	Enabled bool                     `hcl:"enabled"`
	Config  *TelegramConf            `hcl:"config,block"`
	With    *AlertConfigTelegramWith `hcl:"extends,block"`
}

type TelegramConf struct {
	Token  string `hcl:"token"`
	ChatID string `hcl:"chat_id"`
	// Message *string `hcl:"message"`
}

type AlertConfigTelegramWith struct {
	Extends string  `hcl:"extends,label"`
	ChatID  *string `hcl:"chat_id"`
	// Message *string `hcl:"message"`
}

type MonitorConfig struct {
	Name     string           `hcl:"name,label"`
	Services []*ServiceConfig `hcl:"service,block"`
}

type ServiceConfig struct {
	Type      string        `hcl:"type,label"`
	Upstream  string        `hcl:"upstream"`
	Interval  *int          `hcl:"interval"`
	Threshold *int          `hcl:"threshold"`
	Timeout   *int          `hcl:"timeout"`
	Headers   []*HTTPHeader `hcl:"add_header,block"`
	Alert     *AlertConfig  `hcl:"alert,block"`
}

type HTTPHeader struct {
	Name  string `hcl:"name,label"`
	Value string `hcl:"value"`
}

type Config struct {
	Global  GlobalConfig    `hcl:"global,block"`
	Monitor []MonitorConfig `hcl:"monitor,block"`
}

// LoadConfig loads the configuration from the given file path.
func LoadConfig(path string) (*Config, error) {
	var config Config
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("config file not found: %s", path)
	}

	if stat.IsDir() {
		return nil, fmt.Errorf("config file is a directory: %s", path)
	}

	err = hclsimple.DecodeFile(path, nil, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &config, nil
}

const (
	LogModeAppend LogMode = "append"
)
