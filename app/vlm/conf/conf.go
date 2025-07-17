package conf

import (
	"gopkg.in/yaml.v2"
	"os"
)

type MongoDb struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Rs       string `yaml:"rs"`
	Db       string `yaml:"db"`
	Table    string `yaml:"table"`
}

type Config struct {
	MongoDbConfigs struct {
		DB MongoDb `yaml:"test_db"`
	} `yaml:"mongodb"`
	ImagePath string `yaml:"image_path"`
	WsAddr    string `yaml:"ws_addr""`

	In           int32
	Out          int32
	WriteDBCount int32
	ClientNum    int   `yaml:"client_num"`
	Interval     int64 `yaml:"interval"`
}

var GConfig *Config

func Init(cfgPath string) error {
	yamlFile, err := os.ReadFile(cfgPath)
	if err != nil {
		return err
	}
	c := &Config{}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return err
	}

	GConfig = c
	return nil
}
