package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	ScrapeInterval int               `mapstructure:"scrape_interval"`
	Sources        map[string]Source `mapstructure:"sources"`
}

type Source struct {
	TitleKey       string `mapstructure:"title_key"`
	ImageKey       string `mapstructure:"image_key"`
	Url            string `mapstructure:"url"`
	PagePath       string `mapstructure:"page_path"`
	ScrapeInterval int    `mapstructure:"scrape_interval"`
}

func GetSources() map[string]Source {
	var config Config
	err := viper.Unmarshal(&config)
	if err != nil {
		panic("Failed to parse config..")
	}

	return config.Sources
}

func GetScrapeInterval() int64 {
	return viper.GetInt64("scrape_interval")
}
