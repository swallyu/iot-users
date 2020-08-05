package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/swallyu/iot-users/log"
)

// Config config value
type Config struct {
	Database *Database
}

// Database type
type Database struct {
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	User   string `toml:"user"`
	Pwd    string `toml:"pwd"`
	DbName string `toml:"name"`
}

var Conf Config = Config{}

func LoadConfig() error {
	p, _ := os.Executable()
	fmt.Println(p)
	_, err := toml.DecodeFile("config.toml", &Conf)

	if err != nil {
		log.Println(err)
	}
	return err
}
