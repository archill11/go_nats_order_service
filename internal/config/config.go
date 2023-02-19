package config

import (
	"log"
	"os"
	"sync"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database struct {
		Host       string `yaml:"host"`
		Port       string `yaml:"port"`
		Username   string `yaml:"user"`
		Password   string `yaml:"pass"`
		DBname     string `yaml:"dbname"`
		DriverName string `yaml:"driverName"`
	} `yaml:"database"`

	StanListener struct {
		ClusterName string `yaml:"cluster_name"`
		ClientId    string `yaml:"client_id"`
		DurableName string `yaml:"durable_name"`
		Subject     string `yaml:"subject"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
	} `yaml:"stan-listener"`

	HttpServer struct {
		Port string `yaml:"port"`
	} `yaml:"http-server"`
}

var (
	config Config
	once   sync.Once
)

// InitFile Метод инициализирует значения конфигурации из файла config.yaml
func Get() Config {
	once.Do(func() {
		file, err := os.Open("../../config.yaml")
		if err != nil {
			log.Println(err)
		}
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		err = decoder.Decode(&config)
		if err != nil {
			log.Println(err)
		}
	})
	return config
}
