package setting

import (
	"encoding/json"
	"log"
	"os"
)

type ConfigDatabase struct {
	Host string `json:"host"`
	Port int `json:"port"`
	User string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type ConfigServer struct {
	Mode int `json:"mode"`
	Port int `json:"port"`
}

type Config struct {
	DB ConfigDatabase `json:"db"`
	Server ConfigServer `json:"server"`
}

var (
	c *Config
	cfgPath string
)

func Setup(configPath string) {
	cfgPath = configPath
	_ = GetGlobalConfig()
}

func LoadConfig(path string) Config{
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("[loadConfig]: %s\n", err)
	}
	defer file.Close()
	newCfg := Config{}
	err = json.NewDecoder(file).Decode(&newCfg)
	if err != nil {
		log.Fatalf("[loadConfig Decode]: %s\n", err)
	}
	return newCfg
}

func GetGlobalConfig() Config {
	if c == nil {
		newCfg := LoadConfig(cfgPath)
		c = &newCfg
	}
	return *c
}