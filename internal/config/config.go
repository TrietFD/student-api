package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
 	Addr string `yaml:"address" env-required:"true"`
}

// env-default:"production"

type Config struct {
	Env         string `yaml:"env" env-default:"local"`
	StoragePath string `yaml:"storage_path"`
	
	HTTPServer struct {
		Address string `yaml:"address"`
	} `yaml:"http_server"`

	// MySQL specific configuration
	DBHost     string `yaml:"db_host"`
	DBPort     int    `yaml:"db_port"`
	DBUser     string `yaml:"db_user"`
	DBPassword string `yaml:"db_password"`
	DBName     string `yaml:"db_name"`
}

func MustLoad() *Config{
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	
	if configPath == "" {
		flags := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)	
	}

	var cfg Config


	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Can not read ocnffig file: %s", err.Error())
	}

	return &cfg
}