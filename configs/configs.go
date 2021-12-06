package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DbConfig       DatabaseConfig    `mapstructure:"databases_configs"`
	ScrapeInterval int               `mapstructure:"scrape_interval"`
	Sources        map[string]Source `mapstructure:"sources"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Database string `mapstructure:"database"`
	Password string `mapstructure:"password"`
}

type Source struct {
	TitleKey       string `mapstructure:"title_key"`
	ImageKey       string `mapstructure:"image_key"`
	Url            string `mapstructure:"url"`
	PagePath       string `mapstructure:"page_path"`
	ScrapeInterval int    `mapstructure:"scrape_interval"`
}

func readConfigs() Config {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic("Failed to parse config..")
	}

	return config
}

func GetDbConfigs() DatabaseConfig {
	return readConfigs().DbConfig
}

func GetSources() map[string]Source {
	return readConfigs().Sources
}

func GetScrapeInterval() int64 {
	return viper.GetInt64("scrape_interval")
}
