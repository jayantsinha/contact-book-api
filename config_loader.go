package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Config file struct
type Config struct {
	DB struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Name     string `json:"name"`
	} `json:"db"`
	App struct {
		Domain      string `json:"domain"`
		ListenPort  string `json:"listen-port"`
		Logging     bool   `json:"logging"`
		Environment string `json:"environment"`
	} `json:"app"`
	Logging struct {
		Path string `json:"path"`
	} `json:"logging"`
}

// LoadConfiguration loads the configuration from the config.json file into the app context
func LoadConfiguration() *Config {
	var config Config
	configFile, err := os.Open(preferredConfigFile())
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config
}

// preferredConfigFile loads the config file from the location by order of preference
// Order of preference is
// 1. <current_dir>/config/config.json
// 2. config/config.json
func preferredConfigFile() string {
	var configfile = "config/config.json"
	prefloc, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("Unable to get location of binary. Err: ", err)
		// log.Pringln("Config file missing at Loading configuration from the build itself")
	}
	preferredConfigLoc := prefloc + PATH_SEPARATOR + "config" + PATH_SEPARATOR + "config.json"
	if _, err := os.Stat(preferredConfigLoc); os.IsNotExist(err) {
		configfile = preferredConfigLoc
	}
	return configfile
}
