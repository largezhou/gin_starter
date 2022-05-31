package viper

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"regexp"
)

func init() {
	fmt.Println("viper init start")

	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	configPath := "./config.yaml"
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	configBytes = handleYamlValue(configBytes)

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(configBytes)); err != nil {
		panic(err)
	}

	fmt.Println("viper init end")
}

// 处理 yaml 配置中的环境变量和默认值
func handleYamlValue(configBytes []byte) []byte {
	re := regexp.MustCompile("\\$\\{([^:]*):([^}]*)}")
	configBytes = re.ReplaceAllFunc(configBytes, func(bytes []byte) []byte {
		findBytes := re.FindSubmatch(bytes)
		if len(findBytes) != 3 {
			return bytes
		}
		envName := string(findBytes[1])
		value := viper.Get(envName)
		if value == nil {
			return findBytes[2]
		} else {
			return []byte(fmt.Sprintf("%v", value))
		}
	})

	return configBytes
}
