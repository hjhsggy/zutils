package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Conf define
type Conf struct {
	Env string    `toml:"env"`
	App appConfig `toml:"app"`
}

// appConfig struct
type appConfig struct {
	Name    string `toml:"name"`
	LogPath string `toml:"log_path"`
	LogName string `toml:"log_name"`
}

// App config
var App = &appConfig{}

// Env string
var Env string

func init() {
	filePath, err := filepath.Abs("./config/config.toml")
	if err != nil {
		filePath, err = filepath.Abs("../config/config.toml")
		if err != nil {
			panic(err)
		}
	}

	var conf Conf
	if _, err := toml.DecodeFile(filePath, &conf); err != nil {
		panic(err)
	}

	Env = conf.Env
	App = &conf.App
}
