package config

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"regexp"
)

var Config struct {
	App struct { // 应用基础配置
		Host      string  // 监听 IP
		Port      string  // 监听 端口
		Env       string  // 环境
		Debug     bool    // 是否开启 debug
		Timezone  string  // 时区
		Key       string  // 加密密钥
		DistRange float64 // 定位范围
	}
	Log struct { // 日志配置
		Level  string // 日志级别
		Stdout bool   // 是否同时输出到终端
	}
	Mysql struct { // mysql 配置
		Dsn string // 连接
	}
	Redis struct { // redis 配置
		Host     string
		Password string
		Db       int
	}
}

func init() {
	initViper()

	if err := viper.Unmarshal(&Config); err != nil {
		panic(err)
	}
}

func initViper() {
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
