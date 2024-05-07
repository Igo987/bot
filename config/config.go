package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

var Cfg *Config

var ErrConfig = errors.New("fatal error config file")
var ErrUnmarshal = errors.New("unable to unmarshal config file")

func init() {
	var err error
	Cfg, err = LoadConfig()

	if err != nil {
		panic(err)
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath("../../config")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, ErrConfig
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, ErrUnmarshal
	}

	return &cfg, nil
}
func ReadConfig() (*Config, error) {
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, ErrConfig
	}

	var cfg Config

	err = viper.Unmarshal(&cfg)

	if err != nil {
		return nil, ErrUnmarshal
	}

	return &cfg, err
}

type dataBase struct {
	DBHost     string `yaml:"dbHost"`
	DBPort     string `yaml:"dbPort"`
	DBName     string `yaml:"dbName"`
	DBUser     string `yaml:"dbUser"`
	DBPassword string `yaml:"dbPassword"`
}

func (d *dataBase) String() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", d.DBUser, d.DBPassword, d.DBHost, d.DBPort, d.DBName)
}

type API struct {
	APIKey string `yaml:"apiKey"`
	URL    string `yaml:"url"`
}

type Port struct {
	Port string `yaml:"port"`
}

type Telegram struct {
	Token string `yaml:"token"`
}

type Config struct {
	DataBase dataBase `yaml:"database"`
	API      API      `yaml:"api"`
	Telegram Telegram `yaml:"telegram"`
	Port     Port     `yaml:"port"`
	Waiting  int      `yaml:"waitingTime"`
}

func (c *Config) GetToken() string {
	return c.Telegram.Token
}

func (c *Config) GetAPIKey() string {
	return c.API.APIKey
}

func (c *Config) GetPort() string {
	return c.Port.Port
}

func (c *Config) GetConnString() string {
	return c.DataBase.String()
}

func (c *Config) GetAPIURL() string {
	return c.API.URL
}
