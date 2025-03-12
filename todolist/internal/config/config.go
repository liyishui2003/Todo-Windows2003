package config

import (
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

type Config struct {
	JWTSecret string `mapstructure:"JWT_SECRET"`  // JWT 密钥
	User      string `mapstructure:"DB_USER"`     // 数据库用户名
	Password  string `mapstructure:"DB_PASSWORD"` // 数据库密码
	Host      string `mapstructure:"DB_HOST"`     // 数据库主机
	Port      string `mapstructure:"DB_PORT"`     // 数据库端口
	Name      string `mapstructure:"DB_NAME"`     // 数据库名称
}

var Cfg Config // 定义全局的配置对象

func LoadConfig() (config Config, err error) {
	/*
		    viper.SetDefault("JWT_SECRET", "eat-pray-love-0120") // 默认 JWT 密钥
			viper.SetDefault("DB_USER", "root")                  // 默认数据库用户名
			viper.SetDefault("DB_PASSWORD", "010101")            // 默认数据库密码
			viper.SetDefault("DB_HOST", "localhost")             // 默认数据库主机
			viper.SetDefault("DB_PORT", "3306")                  // 默认数据库端口
			viper.SetDefault("DB_NAME", "todolist")              // 默认数据库名称
			viper.AutomaticEnv()
	*/

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("无法获取当前文件路径")
	}
	dir := filepath.Dir(filename) // 获取当前文件所在的目录

	// 设置配置文件名和路径
	viper.SetConfigName("config") // 配置文件名称（不带扩展名）
	viper.SetConfigType("yaml")   // 配置文件类型
	viper.AddConfigPath(dir)      // 配置文件路径（当前目录）

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return
}
