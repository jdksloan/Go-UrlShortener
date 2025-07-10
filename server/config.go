package server

import (
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Port        string `koanf:"port"`
	ApiPrefix   string `koanf:"api_prefix"`
	ApiVersion  int    `koanf:"api_version"`
	RedirectUrl string `koanf:"redirect_url"`
	LogLevel    zerolog.Level
}

func LoadConfig(filePath string) (*Config, error) {
	result := &Config{LogLevel: zerolog.InfoLevel}
	k := koanf.New(".")
	err := k.Load(file.Provider(filePath), dotenv.Parser())
	if err != nil {
		return nil, err
	}
	err = k.Unmarshal("", &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
