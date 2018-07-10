package config

import (
	"fmt"
	"log"
	"reflect"

	"github.com/spf13/viper"
)

type Config struct {
	BtcConf BtcConf
}

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
}
func LoadConfig() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return nil, err
	}
	cfg := &Config{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
		return nil, err
	}
	return cfg, nil
}
func CheckConfig(c *Config, cnames []string) error {
	s := reflect.ValueOf(c).Elem()
	for _, v := range cnames {
		val := s.FieldByName(v)
		typ := val.Type()
		def := reflect.New(typ).Elem()
		flag := 0
		i := 0
		for ; i < val.NumField(); i++ {
			if val.Field(i).Interface() == def.Field(i).Interface() {
				flag++
			}
		}
		if i == flag {
			return fmt.Errorf("%v is not find in config", v)
		}
	}
	return nil
}

type BtcConf struct {
	IP     string //地址
	Port   string //端口号
	User   string //账户
	Passwd string //密码
}
