package config

import (
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type CloudSecret struct {
	Key    string `destructure:"key"`
	Secret string `destructure:"secret"`
}
type Config struct {
	Addr    string      `destructure:"addr"`
	Port    int         `destructure:"port"`
	Token   string      `destructure:"token"`
	Ali     CloudSecret `destructure:"ali"`
	Tencent CloudSecret `destructure:"tencent"`
	Version bool        `destructure:"version"`
}

var Cfg Config

func bindFlag(key string, defaultValue interface{}, usage string) {
	// 设置默认值
	viper.SetDefault(key, defaultValue)

	// 根据 defaultValue 的类型选择合适的 flag 绑定
	switch v := defaultValue.(type) {
	case bool:
		pflag.Bool(key, v, usage)
	case int:
		pflag.Int(key, v, usage)
	case string:
		pflag.String(key, v, usage)
	case float64:
		pflag.Float64(key, v, usage)
	default:
		// 如果不常用类型，可以直接转 string
		pflag.String(key, "", usage)
	}
	// viper 绑定 pflag
	_ = viper.BindPFlag(key, pflag.Lookup(key))
}

func init() {
	// 1. 默认值
	bindFlag("addr", "", "Server Address")
	bindFlag("port", 8080, "Server Port")
	bindFlag("token", "", "Application Token")
	bindFlag("version", false, "Print version and exit")
	bindFlag("ali.key", "", "Ali AccessKey Id")
	bindFlag("ali.secret", "", "Ali AccessKey Secret")
	bindFlag("tencent.key", "", "Tencent AccessKey Id")
	bindFlag("tencent.secret", "", "Tencent AccessKey Secret")

	// 2. 配置文件 (支持 yaml/json/toml)
	viper.SetConfigName("config") // 文件名（不带扩展名）
	viper.AddConfigPath(".")      // 查找路径
	viper.SetConfigType("yaml")   // 也可 json/toml

	_ = viper.ReadInConfig() // 读取配置文件

	// 3. 环境变量（支持 APP_ADDR / APP_PORT / APP_DEBUG）
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("FW")
	viper.AutomaticEnv()

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)

	// 5. 绑定到全局变量
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}
}
