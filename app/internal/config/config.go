package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"sync"
)

type Config struct {
	HTTP struct {
		IP   string `yaml:"ip" env:"HTTP-IP"`
		Port int    `yaml:"port" env:"HTTP-PORT"`
	} `yaml:"http"`
	Redis struct {
		Host     string `yaml:"host" env:"REDIS_HOST" env-required:"true"`
		Port     string `yaml:"port" env:"REDIS_PORT" env-required:"true"`
		Database string `yaml:"database" env:"REDIS_DATABASE" env-required:"true" env-default:"0"`
	} `yaml:"redis"`
}

const (
	EnvConfigPathName  = "CONFIG-PATH"
	FlagConfigPathName = "config"
)

var configPath string
var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		flag.StringVar(&configPath, FlagConfigPathName, "configs/config.local.yaml", "this is app config file")
		flag.Parse()

		log.Print("config init")

		if configPath == "" {
			configPath = os.Getenv(EnvConfigPathName)
		}

		if configPath == "" {
			log.Fatal("config path is required")
		}

		instance = &Config{}

		if err := cleanenv.ReadConfig(configPath, instance); err != nil {
			log.Println("Cant`t read environment variables from neither .yaml nor .env")
			log.Println(err)

			err := cleanenv.ReadEnv(instance)
			if err != nil {
				help, _ := cleanenv.GetDescription(instance, nil)
				log.Println(help)
				log.Fatalln(err)
			}
		}
	})
	return instance
}
