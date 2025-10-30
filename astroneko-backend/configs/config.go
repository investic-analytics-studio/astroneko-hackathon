package configs

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config struct
type Config struct {
	App         `mapstructure:"app"`
	Postgres    `mapstructure:"postgres"`
	Firebase    `mapstructure:"firebase"`
	ExternalURL `mapstructure:"external_url"`
}

// App struct
type App struct {
	Debug   bool   `mapstructure:"debug"`
	Env     string `mapstructure:"env"`
	Port    string `mapstructure:"port"`
	JWT     string `mapstructure:"jwt"`
	Project string `mapstructure:"project"`
}

// Postgres struct
type Postgres struct {
	InstanceConnectionName string `mapstructure:"instance_connection_name"`
	Host                   string `mapstructure:"host"`
	Port                   string `mapstructure:"port"`
	Username               string `mapstructure:"username"`
	Password               string `mapstructure:"password"`
	DbName                 string `mapstructure:"database"`
	SSLMode                bool   `mapstructure:"sslmode"`
}

type Firebase struct {
	Credential string `mapstructure:"credential"`
	WebAPIKey  string `mapstructure:"web_api_key"`
}

type ExternalURL struct {
	AstronekoURL string `mapstructure:"astroneko_url"`
	Token        string `mapstructure:"token"`
}

var config Config

// InitViper func
func InitViper(path string) {
	// Load prod-config.yml
	if err := getConfig(path, "prod-config"); err != nil {
		panic(err)
	}
}

// GetViper func
func GetViper() *Config {
	return &config
}

func getConfig(path, configName string) error {
	viper.SetConfigName(configName)
	viper.SetConfigType("yml")
	viper.AddConfigPath(path)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file has changed: ", e.Name)
	})

	return viper.Unmarshal(&config)
}
