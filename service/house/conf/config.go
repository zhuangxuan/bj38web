package conf

import (
	"os"

	"github.com/spf13/viper"
)

func InitConfig() error {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/conf")
	return viper.ReadInConfig()
}
