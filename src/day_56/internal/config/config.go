package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

func Load(cfgFile, env string) error {
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./configs/")
		viper.SetConfigName("config." + env)
	}
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置失败:%w", err)
	}
	log.Println("加载配置文件：", viper.ConfigFileUsed())
	return nil
}
