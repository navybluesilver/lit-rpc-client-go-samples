package config

import (
	"fmt"
	viper "github.com/spf13/viper"
)

func GetString(config string) (s string) {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath("/home/navybluesilver/.navybluesilver")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return viper.GetString(config)
}
