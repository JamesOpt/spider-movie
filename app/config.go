package app

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
)

var Config *viper.Viper

func init()  {

	Config = viper.New()
	Config.AddConfigPath(filepath.Join(GetRootPath(nil), "config"))
	Config.SetConfigType("yaml")
	Config.SetConfigName("app")

	//尝试进行配置读取
	if readFileErr := Config.ReadInConfig(); readFileErr != nil {
		log.Panic(readFileErr)
	}
}

/**
	获取目录
 */
func GetRootPath(dirPath interface{}) string {
	rootPath, _ := os.Getwd()
	if dirPath == nil {
		return rootPath
	}

	return filepath.Join(rootPath, dirPath.(string))
}

