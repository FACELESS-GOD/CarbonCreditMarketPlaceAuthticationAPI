package Configurator

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDRIVER      string `mapstructure:"DBDRIVER"`
	DBCONNSTRING  string `mapstructure:"DBCONNSTRING"`
	RDBCONNSTRING string `mapstructure:"RDBCONNSTRING"`
	ADDRESS       string `mapstructure:"ADDRESS"`
	JwtSecretKey  string `mapstructure:"JWTKEY"`
}

func NewConfigurator(Path string) (Config, error) {
	conf := Config{}

	if len(Path) < 1 {
		return conf, errors.New("Empty Path.")
	}

	viper.AddConfigPath(Path)

	viper.SetConfigName("app")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
		return conf, err
	}

	err = viper.Unmarshal(&conf)

	if err != nil {
		log.Fatal(err)
		return conf, err
	}

	return conf, nil

}
