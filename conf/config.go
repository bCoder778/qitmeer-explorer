package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/bCoder778/log"
	"os"
	"sync"
)

const (
	configFile = "config.toml"
)

var Setting *Config
var once sync.Once

func init() {

	once.Do(func() {
		if Exist(configFile) {
			if _, err := toml.DecodeFile(configFile, &Setting); err != nil {
				fmt.Println(err)
			}
		}
	})
}

type Config struct {
	Api       *Api       `toml:"api"`
	Rpc       *Rpc       `toml:"rpc"`
	DB        *DB        `toml:"db"`
	Log       *Log       `toml:"log"`
	Resources *Resources `toml:"resources"`
}

type Api struct {
	Port int `toml:"port"`
}

type Rpc struct {
	Host     string `toml:"host"`
	Admin    string `toml:"admin"`
	Password string `toml:"password"`
}

type DB struct {
	DBType   string `toml:"dbtype"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Address  string `toml:"address"`
	DBName   string `toml:"dbname"`
}

type Resources struct {
	CPUNumber int `toml:"cpu-number"`
	GCPercent int `toml:"gc-percent"`
}

type Log struct {
	Mode  log.Mode  `toml:"mode"`
	Level log.Level `toml:"level"`
}

func Exist(fileName string) bool {
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
