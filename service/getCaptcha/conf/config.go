package conf

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func InitConfig() error {
	workDir, _ := os.Getwd()
	fmt.Println("配置目录", workDir)
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/conf")
	return viper.ReadInConfig()
}
