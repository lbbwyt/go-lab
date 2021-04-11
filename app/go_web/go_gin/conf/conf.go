package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Web struct {
	Port string `yaml:"port"`
}

type Kafka struct {
	Addrs   []string `yaml:"addrs"`
	Version string   `yaml:"version"`
	Group   string   `yaml:"group"`
}

type MongoDb struct {
	Addr     string `yaml:"addr"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Rs       string `yaml:"rs"`
	Db       string `yaml:"db"`
}

type Logger struct {
	Debug        bool `yaml:"debug"`
	ReportCaller bool `yaml:"report_caller"`
}

type ClickHouse struct {
	Addrs       []string `yaml:"addrs"`
	User        string   `yaml:"user"`
	Password    string   `yaml:"password"`
	StatisticDb string   `yaml:"statistic_db"`
}

type Config struct {
	MongoDbConfigs struct {
		DB MongoDb `yaml:"test_db"`
	} `yaml:"mongodb"`
	WebServer  Web        `yaml:"web"`
	Kafka      Kafka      `yaml:"kafka"`
	Logger     Logger     `yaml:"logger"`
	ClickHouse ClickHouse `yaml:"click_house"`
}

var GConfig *Config

func Init(cfgPath string) error {
	yamlFile, err := ioutil.ReadFile(cfgPath)
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
